package ws

var _wsManager WsManager

type WsManager struct {
	wsPool []WsConn
}
