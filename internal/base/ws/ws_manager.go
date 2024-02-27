package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"wstester/internal/base/code"
	"wstester/pkg/id"
	"wstester/pkg/log"

	"github.com/pkg/errors"
)

type WsManager struct {
	clientPoolLocker sync.Mutex
	clientPool       map[int64]*WsClient

	resultChStack      []chan []byte
	messageResultChMap map[int]chan []byte
	notifyListenCh     chan *WsClient

	//id generator
	serverIDGenerator id.IDGenerator
	msgIDGenerator    id.IDGenerator

	//----------------------------------------
	serverPool      map[int64]*WsServer
	serverPoolLoker sync.Mutex

	serverClientMap map[int64]int64

	serverResp     map[int64]chan []byte
	serverRespLock sync.Mutex

	serverMessage     map[int64]int64
	serverMessageLock sync.Mutex
}

func NewWsManager() *WsManager {
	return &WsManager{
		clientPool:         make(map[int64]*WsClient),
		resultChStack:      make([]chan []byte, 0),
		messageResultChMap: make(map[int]chan []byte),
		notifyListenCh:     make(chan *WsClient),

		serverPool:    make(map[int64]*WsServer),
		serverResp:    make(map[int64]chan []byte),
		serverMessage: make(map[int64]int64),
	}
}

func (m *WsManager) InitAndRegisterClient(platformID int64, wsUrl string) error {
	wsCli := NewWsClient(wsUrl)
	if err := wsCli.ConnectToServer(); err != nil {
		return err
	}

	m.clientPoolLocker.Lock()
	defer m.clientPoolLocker.Unlock()
	m.clientPool[platformID] = wsCli
	return nil
}

func (m *WsManager) FetchPlatformClient(platformID int64) (*WsClient, bool) {
	wsCli, exist := m.clientPool[platformID]
	return wsCli, exist
}

func (m *WsManager) SendMessageToPlatform(platformID int64, msg []byte) (resultCh chan []byte, err error) {
	var msgID int
	wsCli, exist := m.clientPool[platformID]
	if !exist {
		err = errors.New("login first")
		return
	}

	resultCh = m.fetchResultChannel()
	msgID, err = m.resolveMessageID(msg)
	if err != nil {
		return
	}
	m.rememberMessageResultCh(msgID, resultCh)

	err = wsCli.WriteMessage(msg)
	if err != nil {
		return
	}

	return
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

func (m *WsManager) rememberMessageResultCh(msgID int, ch chan []byte) {
	m.messageResultChMap[msgID] = ch
}

func (m *WsManager) fetchResultChannel() chan []byte {
	for len(m.resultChStack) > 0 {
		return m.resultChStack[len(m.resultChStack)-1]
	}

	return make(chan []byte)
}

func (m *WsManager) createWsServer() *WsServer {
	id := m.serverIDGenerator.GetID()
	ws := newWsServer(id)

	m.serverPoolLoker.Lock()
	defer m.serverPoolLoker.Unlock()

	m.serverPool[id] = ws
	return ws
}

func (m *WsManager) UpgradeHttpToWsAndServer(w http.ResponseWriter, r *http.Request) error {
	ws := m.createWsServer()

	if err := ws.node.UpgradeHttp(w, r); err != nil {
		log.Errorf("upgrade http to websocket error: %s", err.Error())
		return err
	}

	go func() {
		ws.Serve()
	}()
	return nil
}

func (m *WsManager) SendMessage(serverID int64, rawMessage string) error {
	message := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rawMessage), &message); err != nil {
		log.Errorf("unmarshal message: %s error: %s", rawMessage, err.Error())
		return code.ERR_JSON_ERROR
	}

	id := m.msgIDGenerator.GetID()
	message["id"] = id
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Errorf("marshal message: %v error: %s", message, err.Error())
		return code.ERR_JSON_ERROR
	}
	m.serverMessageLock.Lock()
	m.serverMessage[id] = serverID
	m.serverMessageLock.Unlock()

	clientID, exist := m.serverClientMap[serverID]
	if !exist {
		return code.ERR_NOT_LOGIN
	}

	wc, exist := m.clientPool[clientID]
	if !exist {
		return code.ERR_NOT_LOGIN
	}

	if err := wc.WriteMessage(msgBytes); err != nil {
		return err
	}

	if _, exist := m.serverResp[serverID]; !exist {
		respCh := make(chan []byte)
		m.serverRespLock.Lock()
		m.serverResp[serverID] = respCh
		m.serverRespLock.Unlock()
	}

	return nil
}
