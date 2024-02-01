package video

import (
	"errors"
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
	return nil
}

func (v *VideoService) SendMessage(token string, message interface{}) error {
	value, ok := v.wsPool.Load(token)
	if !ok {
		return errors.New("ws is not init")
	}

	value.()
}
