package message

import (
	"wstester/internal/base/data"
	"wstester/internal/entity"
)

type MessageRepo struct {
	data *data.Data
}

func NewMessageRepo(data *data.Data) *MessageRepo {
	return &MessageRepo{
		data: data,
	}
}

func (m *MessageRepo) Insert(message *entity.Message) error {
	return m.data.DB.Create(message).Error
}

func (m *MessageRepo) QueryMessageList(message *entity.Message, pageIndex, pageSize int) (out []entity.Message, err error) {
	tx := m.data.DB.Model(&entity.Message{})

	if message != nil && len(message.Method) > 0 {

	}
	return
}
