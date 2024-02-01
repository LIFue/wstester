package video

import (
	"errors"
	"sync"
	"wstester/internal/base/wsutil"
)

type VideoService struct {
	wsPool *sync.Map
}

func NewVideoService() *VideoService {
	return &VideoService{
		wsPool: new(sync.Map),
	}
}

func (v *VideoService) RegisterWsUrl(token, wsUrl string) error {
	wu, err := wsutil.NewWsUtil(wsUrl, true)
	if err != nil {
		return err
	}

	v.wsPool.Store(token, wu)
	return nil
}

func (v *VideoService) SendMessage(token string, message interface{}) (string, error) {
	value, ok := v.wsPool.Load(token)
	if !ok {
		return "", errors.New("ws is not init")
	}

	wsUtil, ok := value.(*wsutil.WsUtil)
	if !ok {
		return "", errors.New("wrong ws type")
	}

	if wsUtil == nil {
		return "", errors.New("ws is not est")
	}

	msgID, err := wsUtil.SendMsg(message)
	if err != nil {
		return "", err
	}
	resp := wsUtil.GetResp(msgID)

	return string(resp), nil
}

func (v *VideoService) CloseOldConnection(token string) error {
	wu, err := v.loadWsUtil(token)
	if err != nil {
		return nil
	}
}

func (v *VideoService) loadWsUtil(token string) (*wsutil.WsUtil, error) {
	value, ok := v.wsPool.Load(token)
	if !ok {
		return nil, errors.New("ws is not init")
	}

	wsUtil, ok := value.(*wsutil.WsUtil)
	if !ok {
		return wsUtil, errors.New("wrong ws type")
	}
}
