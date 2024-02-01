package jsonbase

type IJsonBase interface {
	SetJsonID(string)
}

type JsonBase struct {
	JosnID string
}

func (j *JsonBase) SetJsonID(id string) {
	j.JosnID = id
}
