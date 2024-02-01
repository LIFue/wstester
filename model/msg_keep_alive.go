package model

type KeepAlive struct {
	Expires int    `json:"expires"`
	Date    string `json:"data"`
}

func NewKeepAlive(expires int, data string) KeepAlive {
	return KeepAlive{expires, data}
}

type RespKeepAlive struct {
	Result struct {
		Date string `json:"date"`
	} `json:"result"`
	Sid string `json:"sid"`
	ID  int64  `json:"id"`
}
