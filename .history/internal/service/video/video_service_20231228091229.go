package video

import "wstester/internal/base/wsutil"

type VideoService struct {
	wsUrl string
}

func (v *VideoService) SendMessage() {
wu, err := wsutil.NewWsUtil(v.wsUrl, true)

if err != nil {
return	
}

wu.
}
