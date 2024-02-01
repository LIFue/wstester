package ws

import (
	"sync"

	"github.com/pkg/errors"
)

type WsManager struct {
	clientPoolLocker sync.Mutex
	clientPool       map[string]*WsClient
}

func NewWsManager() *WsManager {
	return &WsManager{
		clientPool: make(map[string]*WsClient),
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

func (m *WsManager) SendMessageToPlatform(platformID string, msg []byte) error {
	wsCli, exist := m.clientPool[platformID]
	if !exist {
		return errors.New("login first")
	}
}
