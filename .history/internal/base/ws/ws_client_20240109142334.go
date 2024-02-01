package ws

type WsClient struct {
	node *wsNode
	url string
}

func NewWsClient() *WsClient {
	return &WsClient{}
}

func(client *WsClient) 
