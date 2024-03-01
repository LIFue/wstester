package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wstester/internal/base/code"
	"wstester/pkg/id"
	"wstester/pkg/log"

	"github.com/pkg/errors"
)

type WsManager struct {
	//id generator
	serverIDGenerator id.IDGenerator
	msgIDGenerator    id.IDGenerator

	//----------------------------------------
	// serverPool      map[int64]*WsServer
	// serverPoolLoker sync.Mutex

	serverResp     map[int64]chan []byte
	serverRespLock sync.Mutex
	serverMap      map[int64]*WsServer
	serverMapLock  sync.Mutex

	serverMessage       map[int64]int64
	serverMessageLock   sync.Mutex
	serverLastUseTs     map[int64]int64
	serverLastUseTsLock sync.Mutex

	clientMap            map[string]*WsClient
	clientMapLock        sync.Mutex
	serverClientMap      map[int64]string
	serverClientMapLock  sync.Mutex
	clientLastSendTs     map[string]int64
	clientLastSendTsLock sync.Mutex

	ctx                   context.Context
	managerTaskCtx        context.Context
	managerTaskcancleFunc context.CancelFunc
}

func NewWsManager(ctx context.Context) *WsManager {
	managerTaskCtx, managerTaskcancleFunc := context.WithCancel(ctx)
	m := &WsManager{
		// serverPool:    make(map[int64]*WsServer),
		serverResp:      make(map[int64]chan []byte),
		serverMessage:   make(map[int64]int64),
		serverLastUseTs: make(map[int64]int64, 0),
		serverMap:       make(map[int64]*WsServer),

		clientMap:        make(map[string]*WsClient),
		serverClientMap:  make(map[int64]string),
		clientLastSendTs: make(map[string]int64),

		ctx:                   ctx,
		managerTaskCtx:        managerTaskCtx,
		managerTaskcancleFunc: managerTaskcancleFunc,
	}
	go m.managerTask()
	return m
}

func (m *WsManager) InitAndRegisterClient(serverID int64, platformID string, wsUrl string, isPublic bool) error {
	// var wsCli *WsClient
	// var exist bool
	log.Infof("m.clientMap: %+v", m.clientMap)
	if _, exist := m.clientMap[platformID]; !exist {
		wsCli := NewWsClient(platformID, wsUrl, isPublic)
		if err := wsCli.ConnectToServer(); err != nil {
			return err
		}
		log.Infof("add to map , lock")
		m.clientMapLock.Lock()
		m.clientMap[platformID] = wsCli
		m.clientMapLock.Unlock()
		log.Infof("add to map , unlock")
		go m.readMessage(wsCli)
	}

	m.serverClientMapLock.Lock()
	m.serverClientMap[serverID] = platformID
	m.serverClientMapLock.Unlock()

	return nil
}

func (m *WsManager) managerTask() {
	t := time.NewTicker(1 * time.Minute)
	keepAlive := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-t.C:
			go m.checkClientUsageAndDeleteUnused()
			go m.checkServerUsageAndDeleteUnused()
		case <-keepAlive.C:
			go m.keepAlive()
		case <-m.managerTaskCtx.Done():
			return
		}
	}
}

func (m *WsManager) resolveMessageID(msg []byte) (int, error) {
	temp := struct {
		ID int `json:"id"`
	}{}

	if err := json.Unmarshal(msg, &temp); err != nil {
		log.Errorf("resolve meessage id error: %s", err.Error())
		return 0, err
	}

	if temp.ID == 0 {
		return 0, errors.New("message id is zero")
	}

	return temp.ID, nil
}

func (m *WsManager) UpgradeHttpToWsAndServer(w http.ResponseWriter, r *http.Request) error {
	ws := newWsServer(m.serverIDGenerator.GetID())

	if err := ws.node.UpgradeHttp(w, r); err != nil {
		log.Errorf("upgrade http to websocket error: %s", err.Error())
		return err
	}

	// ws.node.conn.SetPingHandler(func(appData string) error {
	// 	m.serverLastUseTsLock.Lock()
	// 	m.serverLastUseTs[ws.serverID] = time.Now().Unix()
	// 	m.serverLastUseTsLock.Lock()
	// 	return nil
	// })

	go func() {
		ws.Serve()
	}()
	return nil
}

func (m *WsManager) SendMessage(serverID int64, clientID, rawMessage string) (chan []byte, error) {

	log.Infof("send message : %s", rawMessage)
	message := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rawMessage), &message); err != nil {
		log.Errorf("unmarshal message: %s error: %s", rawMessage, err.Error())
		return nil, code.ERR_JSON_ERROR
	}

	id := m.msgIDGenerator.GetID()
	message["id"] = id
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Errorf("marshal message: %v error: %s", message, err.Error())
		return nil, code.ERR_JSON_ERROR
	}
	m.serverMessageLock.Lock()
	m.serverMessage[id] = serverID
	m.serverMessageLock.Unlock()

	if clientID == "" {
		var exist bool
		clientID, exist = m.serverClientMap[serverID]
		if !exist {
			return nil, code.ERR_NOT_LOGIN
		}
	}

	wc, exist := m.clientMap[clientID]
	if !exist {
		return nil, code.ERR_NOT_LOGIN
	}

	if err := wc.WriteMessage(msgBytes); err != nil {
		log.Errorf("client send message error: %s", err.Error())
		return nil, err
	}
	var respCh chan []byte
	if respCh, exist = m.serverResp[serverID]; !exist {
		respCh = make(chan []byte)
		m.serverRespLock.Lock()
		m.serverResp[serverID] = respCh
		m.serverRespLock.Unlock()
	}

	m.clientLastSendTsLock.Lock()
	m.clientLastSendTs[clientID] = time.Now().Unix()
	m.clientLastSendTsLock.Unlock()
	return respCh, nil
}

func (m *WsManager) readMessage(ws *WsClient) {
	for {
		readBytes, err := ws.ReadMessage()
		if err != nil {
			return
		}

		go func() {
			msgID, err := m.resolveMessageID(readBytes)
			if err != nil {
				return
			}

			log.Infof("read msg id: %d", msgID)

			serverID, exist := m.serverMessage[int64(msgID)]
			if !exist {
				return
			}

			log.Infof("serverID: %d", serverID)
			serverCh, exist := m.serverResp[serverID]
			log.Infof("exist: %v serverCh: %+v", exist, serverCh)
			if !exist {
				return
			}
			serverCh <- readBytes
			log.Infof("response: %s", string(readBytes))
		}()
	}
}

func (m *WsManager) checkClientUsageAndDeleteUnused() {
	expireClients := make([]string, 0)
	for clientID, ts := range m.clientLastSendTs {
		if time.Now().Add(-5*time.Minute).Unix() > ts {
			expireClients = append(expireClients, clientID)
			if client, exist := m.clientMap[clientID]; exist {
				client.node.Close()
				m.clientMapLock.Lock()
				delete(m.clientMap, clientID)
				m.clientMapLock.Unlock()
			}
		}
	}
	notUsedClient := make([]string, 0)
	for _, client := range m.clientMap {
		if _, exist := m.clientLastSendTs[client.id]; !exist {
			notUsedClient = append(notUsedClient, client.id)
		}
	}
	for _, clientID := range expireClients {
		delete(m.clientLastSendTs, clientID)
	}
	for _, clientID := range notUsedClient {
		if client, exist := m.clientMap[clientID]; exist {
			client.node.Close()
			m.clientMapLock.Lock()
			delete(m.clientMap, clientID)
			m.clientMapLock.Unlock()
		}
	}
}

func (m *WsManager) checkServerUsageAndDeleteUnused() {
	expireServers := make([]int64, 0)
	for serverID, ts := range m.serverLastUseTs {
		if time.Now().Add(-5*time.Minute).Unix() > ts {
			expireServers = append(expireServers, serverID)

			if server, exist := m.serverMap[serverID]; exist {

				server.node.Close()
				m.serverMapLock.Lock()
				delete(m.serverMap, serverID)
				m.serverMapLock.Unlock()
			}

		}
	}
	notUsedServer := make([]int64, 0)
	for serverID := range m.serverMap {
		if _, exist := m.serverLastUseTs[serverID]; !exist {
			notUsedServer = append(notUsedServer, serverID)
		}
	}
	for _, serverID := range expireServers {
		delete(m.serverLastUseTs, serverID)
	}
	for _, serverID := range notUsedServer {
		if server, exist := m.serverMap[serverID]; exist {
			server.node.Close()
			m.serverMapLock.Lock()
			delete(m.serverMap, serverID)
			m.serverMapLock.Unlock()
		}
	}
}

func (m *WsManager) keepAlive() {
	keepaliveMsg := `{"id": %d,"method":"general.keeplive","params":{"expires":60,"date":"2024-02-26 19:59:33"}}`

	for _, client := range m.clientMap {
		go func(wsClient *WsClient) {
			if wsClient == nil || !wsClient.node.connected {
				log.Errorf("client %s is not connected", wsClient.id)
				m.clientMapLock.Lock()
				delete(m.clientMap, wsClient.id)
				m.clientMapLock.Unlock()
				return
			}
			msgID := m.msgIDGenerator.GetID()
			keepaliveMsg = fmt.Sprintf(keepaliveMsg, msgID)
			log.Infof("try to keep alive client %s, message: %s", wsClient.id, keepaliveMsg)

			respCh, err := m.SendMessage(-1, wsClient.id, keepaliveMsg)
			if err != nil {
				log.Errorf("client: %s keepalive error: %s", wsClient.id, err.Error())
				m.clientMapLock.Lock()
				delete(m.clientMap, wsClient.id)
				m.clientMapLock.Unlock()
				return
			}
			timeout := time.NewTicker(10 * time.Second)
			select {
			case data := <-respCh:
				log.Infof("keepalive response data: %s", string(data))
			case <-timeout.C:
				log.Errorf("keepalive error, response timeout")
				m.clientMapLock.Lock()
				wsClient.node.conn.Close()
				delete(m.clientMap, wsClient.id)
				m.clientMapLock.Unlock()
			}
		}(client)
	}

}
