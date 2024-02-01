package encrypt

import "github.com/google/wire"

var EncryptSet = wire.NewSet(NewSM3, NewMD5, NewEncrypt)
