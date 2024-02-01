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
		m:             sync.Mutex{},
		checkInternal: 1 * time.Minute,
		wsPool:        make(map[string]*WsConn),
	}
}

func (w *WsManager) AddWs(ws *WsConn) {
	w.m.Lock()
	defer w.m.Unlock()

	w.wsPool[ws.id] = ws
}

func AddWs(ws *WsConn) {
	_wsManager.AddWs(ws)
}

func (w *WsManager) ExtendWsLifeTime(wsID string) {
	if condition {

	}
}

func (w *WsManager) Run() {
	ticker := time.NewTicker(w.checkInternal)
	if w.checkInternal == 0 {
		ticker.Stop()
	}
	for range ticker.C {
		now := time.Now()
		for id, ws := range w.wsPool {
			if now.After(ws.liveTime) {
				ws.Close()
				delete(w.wsPool)
			}
		}
	}
}
