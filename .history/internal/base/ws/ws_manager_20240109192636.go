package ws

import (
	"encoding/json"
	"sync"
	"wstester/pkg/log"

	"github.com/pkg/errors"
)

type WsManager struct {
	clientPoolLocker sync.Mutex
	clientPool       map[string]*WsClient

	resultChStack      []chan []byte
	messageResultChMap map[int]chan []byte
}

func NewWsManager() *WsManager {
	return &WsManager{
		clientPool:         make(map[string]*WsClient),
		resultChStack:      make([]chan []byte, 0),
		messageResultChMap: make(map[int]chan []byte),
	}
}

func (m *WsManager) InitAndRegisterClient(platformID string, wsUrl string) error {
	wsCli := NewWsClient(wsUrl)
	if err := wsCli.ConnectToServer(); err != nil {
		return err
	}

	m.clientPoolLocker.Lock()
	defer m.clientPoolLocker.Unlock()
	m.clientPool[platformID] = wsCli
	return nil
}

func (m *WsManager) FetchPlatformClient(platformID string) (*WsClient, bool) {
	wsCli, exist := m.clientPool[platformID]
	return wsCli, exist
}

func (m *WsManager) SendMessageToPlatform(platformID string, msg []byte) (resultCh chan []byte, err error) {
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

func (m *WsManager) fetchResultChannel() (resultCh chan []byte) {
	if len(m.resultChStack) > 0 {
		resultCh = m.resultChStack[len(m.resultChStack)-1]
	} else {
		resultCh  make(chan []byte)
	}

	return resultCh
}
