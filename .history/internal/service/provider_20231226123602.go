package service

import (
	"github.com/google/wire"
	"wstester/internal/service/login"
	"wstester/internal/service/platform"
	"wstester/internal/service/ssh"
)

var ServiceSet = wire.NewSet(login.NewLogin, platform.NewPlatformService, ssh.NewSshService)
