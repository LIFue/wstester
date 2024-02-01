package jsonbase

type IJsonBase interface {
	SetJsonID(int64)
}

type JsonBase struct {
	JosnID int64
}

func (j *JsonBase) SetJsonID(id int64) {
	j.JosnID = id
}
