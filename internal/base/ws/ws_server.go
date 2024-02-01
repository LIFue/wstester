package ws

import (
	"net/rpc"

	"wstester/pkg/jsonrpc"
)

type WsServer struct {
	serverID int64

	node *wsNode
}

func newWsServer(id int64) *WsServer {
	return &WsServer{
		serverID: id,
		node:     NewWsNode(),
	}
}

func (server *WsServer) Serve() {
	rpc.ServeCodec(jsonrpc.NewServerCodec(server.serverID, server.node))
}
