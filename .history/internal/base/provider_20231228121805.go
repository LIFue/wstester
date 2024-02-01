package base

import (
	"wstester/internal/base/data"
	"wstester/internal/base/encrypt"
	"wstester/internal/base/httputil"

	"github.com/google/wire"
)

var BaseSet = wire.NewSet(data.DataSet, encrypt.EncryptSet, httputil.NewHttpUtil())
