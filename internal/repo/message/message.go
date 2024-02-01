package message

import (
	"fmt"
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

func (m *MessageRepo) Insert(message *entity.MessageEntity) error {
	return m.data.DB.Create(message).Error
}

func (m *MessageRepo) QueryMessageList(message *entity.MessageEntity, pageIndex, pageSize int) (messageList []entity.MessageEntity, err error) {
	tx := m.data.DB.Model(&entity.MessageEntity{})

	if message != nil && len(message.Method) > 0 {
		tx = tx.Where("method LIKE ?", fmt.Sprintf("%%%s%%", message.Method))
	}

	err = tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&messageList).Error
	return
}

func (m *MessageRepo) ExistSameMessage(message *entity.MessageEntity) (bool, error) {
	var count int64
	if err := m.data.DB.Model(&entity.MessageEntity{}).Where("method = ?", message.Method).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *MessageRepo) UpdateMessage(message *entity.MessageEntity) error {
	if err := m.data.DB.Model(&entity.MessageEntity{}).Where("method = ?", message.Method).Updates(message).Error; err != nil {
		return err
	}
	return nil
}
