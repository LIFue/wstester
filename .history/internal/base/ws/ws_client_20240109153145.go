package ws

type WsClient struct {
	id   string
	node *wsNode
	url  string
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
