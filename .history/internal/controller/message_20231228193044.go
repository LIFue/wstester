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
	messageList, err := m.messageService.QueryMessageList(&req.Message, req.PageIndex, req.PageSize)
	if err != nil  {
		return err
	}

	resp.MessageList = 
}
