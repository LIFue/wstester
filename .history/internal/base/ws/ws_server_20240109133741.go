package ws

type WsServer struct {
	node *wsNode
}

func NewWsServer() *WsServer {
	return &WsServer{}
}
