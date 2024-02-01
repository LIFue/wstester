package ws

import "sync"

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
	return nil
}
