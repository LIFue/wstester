package schema

import "wstester/internal/entity"

type ReqQueryMessageList struct {
	entity.Message
	PageIndex int `json:"page_index" form:"page_index"`
	PageSize  int `json:"page_size" form:"page_size"`
}

type RespQuerMessageList struct {
	MessageList []entity.Message `json:"message_list" form:"message_list"`
}
