package message

import "wstester/internal/repo/message"

type MessageService struct {
	messageRepo *message.MessageRepo
}

func NewMessageService(messageRepo *message.MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}
