package message

import (
	"wstester/internal/entity"
	"wstester/internal/repo/message"
)

type MessageService struct {
	messageRepo *message.MessageRepo
}

func NewMessageService(messageRepo *message.MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (m *MessageService) QueryMessageList(messageInfo *entity.MessageEntity, pageIndex, pageSize int) ([]entity.MessageEntity, error) {
	return m.messageRepo.QueryMessageList(messageInfo, pageIndex, pageSize)
}

func (m *MessageService) UpdateMessage(messageInfo *entity.MessageEntity) error {
	isExist, err := m.messageRepo.ExistSameMessage(messageInfo)
	if err != nil 
	 {
		return err
	}
}
