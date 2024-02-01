package login

import (
	"fmt"
	"net/http"
	"wstester/internal/base/code"
	"wstester/internal/base/encrypt"
	"wstester/internal/base/httputil"
	"wstester/internal/entity"
	"wstester/internal/schema"
)

type LoginService struct {
	h   *httputil.HttpUtil
	enc *encrypt.Encrypt
}

func NewLogin() *LoginService {
	return new(LoginService)
}

func (l *LoginService) Login(platform *entity.Platform) (string, error) {
	return l.tryLogin(platform)
}

func (l *LoginService) tryLogin(platform *entity.Platform) (string, error) {
	// 1. 获取平台的操作账号
	if !platform.User.IsLegalUser() {
		return "", nil
	}

	loginUrl := platform.GenLoginUrl()
	payload := schema.ReqLoginAuth{
		User: platform.User.Uid,
	}
	//获取登录加密信息
	authRequest := httputil.NewRequest(http.MethodPost, loginUrl, payload)
	respLoginAuth := &schema.RespLoginAuth{}
	if err := l.h.SendRequest(authRequest, respLoginAuth); err != nil {
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
