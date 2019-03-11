package tgo

import (
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"io"
)

type Protocol interface {
	DecodePacket(reader io.Reader) (packets.Packet,error)
	EncodePacket(packet packets.Packet) ([]byte,error)
}
