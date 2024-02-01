package service

import (
	"wstester/internal/service/login"
	"wstester/internal/service/platform"
	"wstester/internal/service/ssh"
	"wstester/internal/service/video"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(login.NewLogin, platform.NewPlatformService, ssh.NewSshService, video.VideoService)
