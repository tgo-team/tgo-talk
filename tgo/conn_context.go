package tgo

import (
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

type ConnContext struct {
	Packet packets.Packet
	Conn Conn
	Server Server
}

func NewConnContext(packet packets.Packet,conn Conn,server Server) *ConnContext {
	return &ConnContext{
		Packet: packet,
		Conn:conn,
		Server: server,
	}
}