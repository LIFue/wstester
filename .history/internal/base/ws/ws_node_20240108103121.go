package ws

import "github.com/gorilla/websocket"

type wsNode struct {
	conn *websocket.Conn
}
