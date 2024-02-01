package controller

import "wstester/internal/service/message"

type MessageController struct {
	messageService *message.MessageService
}
