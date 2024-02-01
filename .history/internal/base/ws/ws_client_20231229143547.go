package ws

import "github.com/gorilla/websocket"

type WsClient struct {
	wsUrl string
	conn  *websocket.Conn
}

func NewWsClient(wsUrl string) *WsClient {
	return &WsClient{
		wsUrl: wsUrl,
	}
}

func (w *WsClient) Conn() error {
	d := websocket.Dialer{}
	d.Dial(w.wsUrl, nil)
}
