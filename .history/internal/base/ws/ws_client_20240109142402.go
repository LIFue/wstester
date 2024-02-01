package ws

type WsClient struct {
	node *wsNode
	url  string
}

func NewWsClient(url string) *WsClient {
	return &WsClient{url: url}
}

func (client *WsClient) ConnectToServer() error {

}
