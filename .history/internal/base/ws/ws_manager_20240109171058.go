package ws

type WsManager struct {
	clientPool map[string]*WsClient
}

func NewWsManager() *WsManager {
	return &WsManager{
		clientPool: make(map[string]*WsClient),
	}
}

func (m *WsManager) InitAndRegisterClient(wsUrl string) {
	wsCli := NewWsClient(wsUrl)
	wsCli.ConnectToServer()
}
