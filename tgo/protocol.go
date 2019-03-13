package tgo

import (
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

type Protocol interface {
	DecodePacket(reader Conn) (packets.Packet,error)
	EncodePacket(packet packets.Packet) ([]byte,error)
}
