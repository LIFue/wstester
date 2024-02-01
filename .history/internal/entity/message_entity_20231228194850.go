package entity

import "gorm.io/gorm"

type MessageEntity struct {
	gorm.Model
	Method  string `gorm:"index" json:"method"`
	Message string `json:"message"`
}

func (m *MessageEntity) TableName() string {
	return "message"
}
