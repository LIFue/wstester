package ws

import "github.com/gorilla/websocket"

type WsClient struct {
	conn *websocket.Conn
}
