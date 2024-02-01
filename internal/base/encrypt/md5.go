package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5 struct {
}

func NewMD5() *MD5 {
	return new(MD5)
}

func (m *MD5) Encode(s string) string {
	h := md5.New()

	h.Write([]byte(s))
	b := h.Sum(nil)
	encodestr := hex.EncodeToString(b)
	return encodestr
}
