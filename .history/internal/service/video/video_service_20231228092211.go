package video

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"wstester/internal/base/wsutil"
	"wstester/pkg/log"

	"github.com/gin-gonic/gin"
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

	wsUtil, ok := value.(*wsutil.WsUtil)
	if !ok {
		return errors.New("wrong ws type")
	}

	if wsUtil == nil {
		return errors.New("ws is not est")
	}

	i, err := wsUtil.SendMsg(message)
	if err != nil {
		return err
	}
	resp := ws.ws.GetResp(msgID)
	m := make(map[string]interface{})
	if err = json.Unmarshal(resp, &m); err != nil {
		log.Errorf("unmarshal resp error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
		return
	}
}
