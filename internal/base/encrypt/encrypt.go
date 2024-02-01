package encrypt

import "wstester/internal/base/code"

type Encrypt struct {
	sm3 *SM3
	md5 *MD5
}

func NewEncrypt(sm3 *SM3, md5 *MD5) *Encrypt {
	return &Encrypt{
		sm3: sm3,
		md5: md5,
	}
}

func (e *Encrypt) Encode(encryptType code.EncryptCode, str string) string {
	switch encryptType {
	case code.EncryptCodeSM3:
		return e.sm3.Encode(str)
	case code.EncryptCodeMD5:
		fallthrough
	default:
		return e.md5.Encode(str)
	}
}
