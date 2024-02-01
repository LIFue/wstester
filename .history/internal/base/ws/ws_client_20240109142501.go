package ws

import "wstester/pkg/log"

type WsClient struct {
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
	if err := client.node.ConnectWsServer(client.url, nil); err != nil {
		log.Errorf("")
	}
}
