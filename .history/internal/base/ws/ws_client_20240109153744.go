package ws

type WsClient struct {
	node *wsNode
	url  string

	buf []byte
}

func NewWsClient(url string) *WsClient {
	return &WsClient{
		node: NewWsNode(),
		url:  url,
	}
}

func (client *WsClient) ConnectToServer() error {
	return client.node.ConnectWsServer(client.url, nil)
}
