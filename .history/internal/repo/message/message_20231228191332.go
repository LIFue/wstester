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

}
