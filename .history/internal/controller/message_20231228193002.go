package controller

import (
	"wstester/internal/schema"
	"wstester/internal/service/message"
)

type MessageController struct {
	messageService *message.MessageService
}

func NewMessageController(messageService *message.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

func (m *MessageController) QueryMessageList(req *schema.ReqQueryMessageList, resp *schema.RespQuerMessageList) error {

}
