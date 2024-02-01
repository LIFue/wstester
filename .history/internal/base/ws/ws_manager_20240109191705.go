package ws

import (
	"sync"

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
		messageResultChMap: make(map[string]map[int]chan []byte),
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
	wsCli, exist := m.clientPool[platformID]
	if !exist {
		err = errors.New("login first")
		return
	}

	resultCh = m.fetchResultChannel()
	m.rememberMessageResultCh(platformID)

	err = wsCli.WriteMessage(msg)
	if err != nil {
		return
	}

	return
}

func (m *WsManager) resolve()

func (m *WsManager) rememberMessageResultCh(platformID string, msgID int, ch chan []byte) {
	platformMessageResultChMap, exist := m.messageResultChMap[platformID]
	if !exist {
		platformMessageResultChMap = make(map[int]chan []byte)
	}

	platformMessageResultChMap[msgID] = ch
	m.messageResultChMap[platformID] = platformMessageResultChMap
}

func (m *WsManager) fetchResultChannel() chan []byte {
	for len(m.resultChStack) > 0 {
		return m.resultChStack[len(m.resultChStack)-1]
	}

	resultCh := make(chan []byte)
	return resultCh
}
