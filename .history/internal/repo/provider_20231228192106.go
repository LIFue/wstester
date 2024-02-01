package repo

import (
	"wstester/internal/repo/message"
	"wstester/internal/repo/platform"
	"wstester/internal/repo/ssh"
	"wstester/internal/repo/user"

	"github.com/google/wire"
)

var RepoSet = wire.NewSet(platform.NewPlatformRepo, ssh.NewSshHostRepo, user.NewUserRepo, message.NewMessageRepo)
