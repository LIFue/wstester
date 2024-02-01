package jsonbase

type IJsonBase interface {
	SetJsonID(string)
}

type JsonBase struct {
	JosnID string
}
