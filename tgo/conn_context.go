package tgo

import (
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

type ConnContext struct {
	Packet packets.Packet
	Conn Conn
}

func NewConnContext(packet packets.Packet,conn Conn) *ConnContext {
	return &ConnContext{
		Packet: packet,
		Conn:conn,
	}
}