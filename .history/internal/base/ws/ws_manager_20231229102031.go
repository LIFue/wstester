package ws

import (
	"sync"
	"time"
)

var _wsManager *WsManager

type WsManager struct {
	m      sync.Mutex
	wsPool map[string]*WsConn

	checkInternal time.Duration
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

	w.wsPool = append(w.wsPool, ws)
}

func AddWs(ws *WsConn) {
	_wsManager.AddWs(ws)
}

func (w *WsManager) Run() {
	ticker := time.NewTicker(w.checkInternal + 1)
	if w.checkInternal == 0 {
		ticker.Stop()
	}
	for {
		select {
		case <-ticker.C:
			case 
		}
	}
}
