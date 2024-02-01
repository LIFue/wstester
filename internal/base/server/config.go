package server

type WebSocket struct {
	Addr string `json:"addr" mapstructure:"addr" yaml:"addr"`
}
