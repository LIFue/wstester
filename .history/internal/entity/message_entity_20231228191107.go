package entity

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Content string `json:"content"`
}
