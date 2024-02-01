package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uid        string `gorm:"index" json:"uid"`
	Username   string
	Password   string
	PlatformID uint
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func (u User) IsLegalUser() bool {
	return len(u.Username) > 0 && len(u.Password) > 0
}
