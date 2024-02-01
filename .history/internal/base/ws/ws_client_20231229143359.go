package ws

import "github.com/gorilla/websocket"

type WsClient struct {
	wsUrl string
	conn  *websocket.Conn
}

func NewWsClient(wsUrl string) *WsClient {
	return
}
