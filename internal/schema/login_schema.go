package schema

type ReqLoginUrlParams struct {
	Time     int64  `json:"time"`
	Zone     int64  `json:"zone"`
	Lang     string `json:"lang"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

type ReqLoginAuth struct {
	User string `json:"user"`
}

type RespLoginAuth struct {
	Host      string `json:"host"`
	Nonce     string `json:"nonce"`
	Algorithm string `json:"algorithm"`
	Realm     string `json:"realm"`
}

type ReqLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type RespLogin struct {
	Ws []string `json:"ws"`
}
