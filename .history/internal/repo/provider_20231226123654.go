package repo

import (
	"github.com/google/wire"
	"wstester/internal/repo/platform"
	"wstester/internal/repo/ssh"
	"wstester/internal/repo/user"
)

var RepoSet = wire.NewSet(platform.NewPlatformRepo, ssh.NewSshHostRepo, user.NewUserRepo)
