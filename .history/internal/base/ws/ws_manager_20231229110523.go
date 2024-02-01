package ws

import (
	"fmt"
	"sync"
	"time"
	"wstester/pkg/log"
)

var _wsManager *WsManager

type WsManager struct {
	m      sync.Mutex
	wsPool map[string]*WsConn

	checkInternal time.Duration
}

func init() {
	fmt.Println("init")
	_wsManager = NewWsManager()
	go _wsManager.Run()
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

func (w *WsManager) Run() {
	ticker := time.NewTicker(w.checkInternal)
	if w.checkInternal == 0 {
		ticker.Stop()
	}
	for range ticker.C {
		now := time.Now()
		for id, ws := range w.wsPool {
			if now.After(ws.liveTime) {
				log.Infof("try close ws conn: %s", id)
				ws.Close()
				delete(w.wsPool, id)
			}
		}
	}
}
