package schema

import "wstester/internal/entity"

type ReqConnectPlatform struct {
	// Ip   string `json:"ip" form:"ip"`
	// Port string `json:"port" form:"port"`
	// entity.Platform
	// User entity.User `json:"user" from:"user"`
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
