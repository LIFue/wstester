package video

import (
	"sync"
	"wstester/internal/base/wsutil"
)

type VideoService struct {
	wsPool sync.Map
}

func (v *VideoService) RegisterWsUrl(token, wsUrl string) error {
	wu, err := wsutil.NewWsUtil(wsUrl, true)
	if err != nil {
		return err
	}

	v.wsPool.Store(token, wu)
}

func (v *VideoService) SendMessage() {
	wu, err := wsutil.NewWsUtil(v.wsUrl, true)

	if err != nil {
		return
	}
}
