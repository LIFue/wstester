package entity

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Method  string `gorm:"index" json:"method"`
	Content string `gorm:"type:size(100)"json:"content"`
}
