package model

import (
	"time"
)

const VERSION = "1.1.0"
const PLATFORM_WEB = "web"

const (
	DEFAULT_PORT     = "3080"
	DEFAULT_PASSWORD = "Hitry@202020"
	DEFAULT_PROTOCOL = "http"
)

type LoginUrlParam struct {
	Time     int64  `json:"time"`
	Zone     int64  `json:"zone"`
	Lang     string `json:"lang"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

func NewLoginUrlParam(version, platform string) LoginUrlParam {
	return LoginUrlParam{
		Time:     time.Now().Unix(),
		Zone:     8,
		Lang:     "cn",
		Version:  version,
		Platform: platform,
	}
}

type LoginPrepareResp struct {
	Host      string `json:"host"`
	Nonce     string `json:"nonce"`
	Algorithm string `json:"algorithm"`
	Realm     string `json:"realm"`
}

type ReqAuthMethodPayload struct {
	User string `json:"user"`
}

func NewReqAuthMethodPayload(uid string) ReqAuthMethodPayload {
	return ReqAuthMethodPayload{
		User: uid,
	}
}

type LoginParam struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func NewLoginParam(user, password string) LoginParam {
	return LoginParam{
		User:     user,
		Password: password,
	}
}

type RespLogin struct {
	Ws []string `json:"ws"`
}

type ReqLoginParam struct {
	Server
	LoginParam
}

func (r *ReqLoginParam) CheckAndFillDefaultValue() error {

	if r.Port == "" {
		r.Port = DEFAULT_PORT
	}

	if r.Password == "" {
		r.Password = DEFAULT_PASSWORD
	}

	if len(r.Protocol) == 0 {
		r.Protocol = DEFAULT_PROTOCOL
	}

	if r.User == "" {
		r.User = "meet"
	}

	return nil
}
