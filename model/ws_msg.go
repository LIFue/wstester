package model

type WsMsg struct {
	ID     int64       `json:"id" validate:"required"`
	Method string      `json:"method" validate:"required"`
	Params interface{} `json:"params"`
}

func NewWsMsg(method string, params interface{}) WsMsg {
	return WsMsg{
		Method: method,
		Params: params,
	}
}
