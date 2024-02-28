package ws

import (
	"context"
	"sync"
	"wstester/pkg/log"
)

type WsClient struct {
	id   string
	node *wsNode
	url  string

	buf []byte

	mu sync.Mutex

	ctx            context.Context
	cancleFuncList []context.CancelFunc

	dataResponse chan []byte
}

func NewWsClient(id string, url string) *WsClient {
	return &WsClient{
		id:             id,
		node:           NewWsNode(),
		url:            url,
		buf:            make([]byte, 8092),
		ctx:            context.Background(),
		cancleFuncList: make([]context.CancelFunc, 0),
	}
}

func (client *WsClient) ConnectToServer() error {
	if err := client.node.ConnectWsServer(client.url, nil); err != nil {
		log.Errorf("connect to server error: %s", err.Error())
		return err
	}
	// ctx, cancleFunc := context.WithCancel(client.ctx)
	// client.cancleFuncList = append(client.cancleFuncList, cancleFunc)
	return nil
}

func (client *WsClient) WriteMessage(msg []byte) error {
	log.Infof("client: %s send message: %s", client.id, string(msg))
	client.mu.Lock()
	defer client.mu.Unlock()
	_, err := client.node.Write(msg)
	return err
}

func (client *WsClient) ReadMessage() ([]byte, error) {
	resp, err := client.node.ReadMessage()
	if err != nil {
		return nil, err
	}
	log.Infof("client resp: %s", string(resp))
	// client.dataResponse <- resp[:n]
	return resp, nil
}
