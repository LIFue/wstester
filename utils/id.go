package utils

import (
	"sync"
	"time"
	"wstester/pkg/log"
)

type IDGenerator struct {
	mu       sync.Mutex
	lastID   int64
	nextTime int64
}

func NewIDGenerator() IDGenerator {
	return IDGenerator{
		nextTime: time.Now().UnixNano(),
	}
}

func (g *IDGenerator) Generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixNano()
	if now < g.nextTime {
		log.Infof("重新生成")
		g.lastID = 0
		g.nextTime = now
	}

	g.lastID++
	return g.lastID
}

func ParseToInt64(id interface{}) int64 {
	switch id.(type) {
	case int64:
		return id.(int64)
	case float64:
		return int64(id.(float64))
	default:
		log.Errorf("type is not supported")
		return 0
	}
}
