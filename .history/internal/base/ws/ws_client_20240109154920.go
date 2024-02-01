package ws

import "sync"

type WsClient struct {
	node *wsNode
	url  string

	buf []byte

	mu sync.Mutex
}

func NewWsClient(url string) *WsClient {
	return &WsClient{
		node: NewWsNode(),
		url:  url,
		buf:  make([]byte, 8092),
	}
}

func (client *WsClient) ConnectToServer() error {
	return client.node.ConnectWsServer(client.url, nil)
}

func (client *WsClient) WriteMessage(msg []byte) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	_, err := client.node.Write(msg)
	return err
}
