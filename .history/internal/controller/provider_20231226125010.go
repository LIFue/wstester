package controller

import "github.com/google/wire"

var ControllerSet = wire.NewSet(NewControllerRegister, NewPlatform, NewArith)
