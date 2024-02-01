package ws

type WsManager struct {
	clientPool map[string]*WsClient
}

func NewWsManager() *WsManager {
	return &WsManager{
		clientPool: map[string]*WsClient{},
	}
}
