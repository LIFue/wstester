package code

import (
	"encoding/json"
)

type ErrorCode int

const (
	ERR_JSON_ERROR ErrorCode = iota + 1
	ERR_NOT_LOGIN
	ERR_TIMEOUT
)

var ErrorCodeMsg = map[ErrorCode]string{
	ERR_JSON_ERROR: "json error",
	ERR_NOT_LOGIN:  "not logind",
	ERR_TIMEOUT:    "timeout",
}

func (e ErrorCode) Error() string {
	tempMsg := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: int(e),
		Msg:  ErrorCodeMsg[e],
	}

	errorMsg, _ := json.Marshal(tempMsg)

	return string(errorMsg)
}
