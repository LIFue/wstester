package model

type Server struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

func NewServer(ip, port, protocol string) Server {
	return Server{
		ip, port, protocol,
	}
}
