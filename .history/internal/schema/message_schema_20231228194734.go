package schema

import "wstester/internal/entity"

type ReqQueryMessageList struct {
	entity.MessageEntity
	PageIndex int `json:"page_index" form:"page_index"`
	PageSize  int `json:"page_size" form:"page_size"`
}

type RespQuerMessageList struct {
	MessageList []entity.MessageEntity `json:"message_list" form:"message_list"`
}
