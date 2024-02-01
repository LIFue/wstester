package ws

import "github.com/gorilla/websocket"

type WsClient struct {
	wsUrl     string
	connected bool
	conn      *websocket.Conn
}

func NewWsClient(wsUrl string) *WsClient {
	return &WsClient{
		wsUrl: wsUrl,
	}
}

func (w *WsClient) Conn() error {
	d := websocket.Dialer{}
	c, _, err := d.Dial(w.wsUrl, nil)
	if err != nil {
		return err
	}
	w.conn = c
	w.connected = true
	return nil
}
