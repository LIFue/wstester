package login

import (
	"errors"
	"fmt"
	"net/http"
	"wstester/internal/base/code"
	"wstester/internal/base/encrypt"
	"wstester/internal/base/httputil"
	"wstester/internal/entity"
	"wstester/internal/schema"
	"wstester/pkg/log"
)

type LoginService struct {
	h   *httputil.HttpUtil
	enc *encrypt.Encrypt
}

func NewLogin(httpUtil *httputil.HttpUtil, encrypt *encrypt.Encrypt) *LoginService {
	return &LoginService{
		h:   httpUtil,
		enc: encrypt,
	}
}

func (l *LoginService) Login(platform *entity.Platform) (string, error) {
	return l.tryLogin(platform)
}

func (l *LoginService) tryLogin(platform *entity.Platform) (string, error) {
	log.Infof("true to login to platform: %+v", platform)
	// 1. 获取平台的操作账号
	if !platform.User.IsLegalUser() {
		log.Errorf("user is not legal")
		return "", errors.New("user is not legal")
	}

	loginUrl := platform.GenLoginUrl()
	payload := schema.ReqLoginAuth{
		User: platform.User.Uid,
	}
	//获取登录加密信息
	authRequest := httputil.NewRequest(http.MethodPost, loginUrl, payload)
	respLoginAuth := &schema.RespLoginAuth{}
	if err := l.h.SendRequest(authRequest, respLoginAuth); err != nil {
		log.Errorf("req login auth error: %s", err.Error())
		return "", err
	}

	// 真实登录，获取websocket链接
	encodePwd := l.enc.Encode(code.EncryptCode(respLoginAuth.Algorithm), platform.User.Password)
	loginPassword := l.enc.Encode(code.EncryptCode(respLoginAuth.Algorithm), fmt.Sprintf("%s:%s:%s:%s", platform.User.Uid, respLoginAuth.Realm, encodePwd, respLoginAuth.Nonce))

	loginPayload := schema.ReqLogin{
		User:     platform.User.Uid,
		Password: loginPassword,
	}
	loginRequest := httputil.NewRequest(http.MethodPost, loginUrl, loginPayload)
	loginResponse := &schema.RespLogin{}
	if err := l.h.SendRequest(loginRequest, loginResponse); err != nil {
		return "", err
	}

	return loginResponse.Ws[0], nil
}
