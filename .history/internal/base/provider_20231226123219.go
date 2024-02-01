package base

import (
	"github.com/google/wire"
	"wstester/internal/base/data"
	"wstester/internal/base/encrypt"
)

var BaseSet = wire.NewSet(data.DataSet, encrypt.EncryptSet)
