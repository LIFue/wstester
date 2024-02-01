package ws

var _wsManager *WsManager

type WsManager struct {
	wsPool []WsConn
}

func init() {
	_wsManager = NewWsManager()
}

func NewWsManager() *WsManager {
	return &WsManager{
		wsPool: make([]WsConn),
	}
}
