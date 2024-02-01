package video

import (
	"sync"
	"wstester/internal/base/wsutil"
)

type VideoService struct {
	wsPool sync.Map
}

func (v *VideoService) Register() {

}

func (v *VideoService) SendMessage() {
	wu, err := wsutil.NewWsUtil(v.wsUrl, true)

	if err != nil {
		return
	}
}
