package ws

type WsManager struct {
	clientPool map[string]*WsClient
}
