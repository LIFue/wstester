package controller

import (
	"wstester/internal/schema"
	"wstester/internal/service/message"
)

type Message struct {
	messageService *message.MessageService
}

func NewMessageController(messageService *message.MessageService, controllerRegister *ControllerRegister) *Message {
	m := &Message{
		messageService: messageService,
	}
	controllerRegister.AddService(m)
	return m
}

func (m *Message) QueryMessageList(req *schema.ReqQueryMessageList, resp *schema.RespQuerMessageList) (err error) {
	resp.MessageList, err = m.messageService.QueryMessageList(&req.MessageEntity, req.PageIndex, req.PageSize)
	return
}
