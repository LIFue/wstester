package ws

type WsClient struct {
	node *wsNode
}

func NewWsClient() *WsClient {
	return &WsClient{}
}
