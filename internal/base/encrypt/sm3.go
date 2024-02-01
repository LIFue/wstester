package encrypt

import (
	"fmt"

	"github.com/tjfoc/gmsm/sm3"
)

type SM3 struct {
}

func NewSM3() *SM3 {
	return new(SM3)
}

func (s *SM3) Encode(str string) string {
	return fmt.Sprintf("%x", sm3.Sm3Sum([]byte(str)))
}
