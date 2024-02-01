package data

import "github.com/google/wire"

var DataSet = wire.NewSet(NewDB, NewData)
