package id

import (
	"math"
	"sync"
	"sync/atomic"
)

type IDGenerator struct {
	idGenerator atomic.Int64
	resetLoker  sync.Mutex
}

func (i *IDGenerator) GetID() int64 {
	tempID := i.idGenerator.Add(1)

	if tempID > math.MaxInt64-10000 {
		i.resetLoker.Lock()
		defer i.resetLoker.Unlock()
		tempID = i.idGenerator.Add(1)
		if tempID > math.MaxInt64-10000 {
			tempID = 1
			i.idGenerator.Store(1)
		}
	}
	return tempID
}
