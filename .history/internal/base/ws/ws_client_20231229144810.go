package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

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

func (w *WsClient) SendAndWait(message string) ([]byte, error)

func (w *WsClient) sendMessage(message string) error {
	if !w.connected {
		return errors.New("ws client is not init")
	}

	return w.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (w *WsClient) receiveMessage() error {
	_, data, err := w.conn.ReadMessage()
	if err != nil {
		return err
	}

}
