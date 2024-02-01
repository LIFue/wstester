package schema

import (
	"wstester/internal/base/jsonbase"
	"wstester/internal/entity"
)

type ReqConnectPlatform struct {
	// Ip   string `json:"ip" form:"ip"`
	// Port string `json:"port" form:"port"`
	// entity.Platform
	// User entity.User `json:"user" from:"user"`
	jsonbase.JsonBase
	entity.Platform
}

type RespConnectPlatform struct {
	Status string
}

type ReqGetPlatformList struct {
	Ip        string `json:"ip" form:"ip"`
	PageIndex int    `json:"page_index" form:"page_index"`
	PageSize  int    `json:"page_size" form:"page_size"`
}

type RespGetPlatformList struct {
	PlatformList []entity.Platform `json:"platform_list"`
}

type ReqSendMessage struct {
	jsonbase.JsonBase
	Method  string `json:"method" form:"method"`
	Message string `json:"message" form:"message"`
}

type RespSendMessage struct {
	Response string `json:"response" form:"response"`
}
