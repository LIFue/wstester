package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uid        string `gorm:"index" json:"uid" form:"uid"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	PlatformID uint   `json:"platform_id" form:"platform_id"`
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
