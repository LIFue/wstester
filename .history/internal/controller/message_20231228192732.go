package controller

import "wstester/internal/service/message"

type MessageController struct {
	messageService *message.MessageService
}

func NewMessageController(messageService *message.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}
