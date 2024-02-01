package entity

import "gorm.io/gorm"

type SshHost struct {
	gorm.Model
	Ip       string
	Port     string
	Username string
	Password string
}

func (s *SshHost) TableName() string {
	return "ssh_host"
}
