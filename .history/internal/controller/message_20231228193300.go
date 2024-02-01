package controller

import (
	"wstester/internal/schema"
	"wstester/internal/service/message"
)

type MessageController struct {
	messageService *message.MessageService
}

func NewMessageController(messageService *message.MessageService, r *ControllerRegister) *MessageController {
	m := &MessageController{
		messageService: messageService,
	}
	r.AddService(m)
	return m
}

func (m *MessageController) QueryMessageList(req *schema.ReqQueryMessageList, resp *schema.RespQuerMessageList) (err error) {
	resp.MessageList, err = m.messageService.QueryMessageList(&req.Message, req.PageIndex, req.PageSize)
	return
}
