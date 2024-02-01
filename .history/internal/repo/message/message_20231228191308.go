package message

import "wstester/internal/base/data"

type MessageRepo struct {
	data *data.Data
}

func NewMessageRepo(data *data.Data) *MessageRepo {
	return &MessageRepo{
		data: data,
	}
}
