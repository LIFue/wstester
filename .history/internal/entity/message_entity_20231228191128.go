package entity

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Method  string `json:"method"`
	Content string `json:"content"`
}
