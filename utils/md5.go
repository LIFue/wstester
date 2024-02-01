package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodeMD5(pass string) string {
	h := md5.New()

	h.Write([]byte(pass))
	b := h.Sum(nil)
	encodestr := hex.EncodeToString(b)
	return encodestr
}
