package ws

import "sync"

var _wsManager *WsManager

type WsManager struct {
	m      sync.Mutex
	wsPool []*WsConn
}

func init() {
	_wsManager = NewWsManager()
}

func NewWsManager() *WsManager {
	return &WsManager{
		wsPool: make([]*WsConn, 0),
	}
}

func (w *WsManager) AddWs(ws *WsConn) {
	w.m.Lock()
	defer w.m.Unlock()

	w.wsPool = append(w.wsPool, *ws)
}
