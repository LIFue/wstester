package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Platform struct {
	gorm.Model
	Ip   string `gorm:"index:idx_platform_ip_port" json:"ip"`
	Port string `gorm:"index:idx_platform_ip_port" json:"port"`
	User User   `json:"user"`
}

func (p *Platform) GenUri() string {
	return fmt.Sprintf("%s://%s:%s", "http", p.Ip, p.Port)
}

func (p *Platform) GenLoginUrl() string {
	return fmt.Sprintf("%s/mesh/user/login?time=%d&zone=8&lang=cn&version=1.1.0&platform=web", p.GenUri(), time.Now().Unix())
}
