package ws

import (
	"context"
	"sync"
	"time"
	"wstester/pkg/log"
)

type WsClient struct {
	node *wsNode
	url  string

	buf []byte

	mu sync.Mutex

	ctx            context.Context
	cancleFuncList []context.CancelFunc

	dataResponse chan []byte
}

func NewWsClient(url string) *WsClient {
	return &WsClient{
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
	ctx, cancleFunc := context.WithCancel(client.ctx)
	client.KeepAlive(ctx)
	client.cancleFuncList = append(client.cancleFuncList, cancleFunc)
	return nil
}

func (client *WsClient) KeepAlive(ctx context.Context) error {
	t := time.NewTicker(50 * time.Second)
	keepaliveMsg := `{"method":"general.keeplive","params":{"expires":60,"date":"2024-02-26 19:59:33"}}`
	for {
		select {
		case <-t.C:
			if err := client.WriteMessage([]byte(keepaliveMsg)); err != nil {
				log.Errorf("keepalive error: send message error: %s", err.Error())
				return err
			}
		case <-ctx.Done():
			log.Info("stop keepAlive")
			return nil
		}
	}
}

func (client *WsClient) SendMessage(msg string) (resp string, err error) {
	if err = client.WriteMessage([]byte(msg)); err != nil {
		return
	}

	return
}

func (client *WsClient) WriteMessage(msg []byte) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	_, err := client.node.Write(msg)
	return err
}

func (client *WsClient) ReadMessage() ([]byte, error) {
	for {
		resp := make([]byte, 8092)
		n, err := client.node.Read(resp)
		if err != nil {
			return nil, err
		}
		client.dataResponse <- resp[:n]
	}
}
