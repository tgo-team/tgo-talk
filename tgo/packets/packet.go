package packets

import (
	"fmt"
	"io"
)

type PacketCodec interface {
	//Decode 解码
	Decode(reader io.Reader) (Packet, error)
	//Encode 编码
	Encode(msg Packet) ([]byte, error)
}

type Packet interface {
	GetFixedHeader() FixedHeader
	String() string
}

type PacketType int
type FixedHeader struct {
	PacketType      PacketType
	Dup             bool
	Qos             byte
	Retain          bool
	RemainingLength int
}
func (fh FixedHeader) String() string {
	return fmt.Sprintf("%s: dup: %t qos: %d retain: %t rLength: %d", PacketNames[uint8(fh.PacketType)], fh.Dup, fh.Qos, fh.Retain, fh.RemainingLength)
}

const (
	None        PacketType = iota
	Connect                = 1
	Connack                = 2
	Message                = 3
	Msgack                 = 4
	Pubrec                 = 5
	Pubrel                 = 6
	Pubcomp                = 7
	Subscribe              = 8
	Suback                 = 9
	Unsubscribe            = 10
	Unsuback               = 11
	Pingreq                = 12
	Pingresp               = 13
	Disconnect             = 14
)

var PacketNames = map[uint8]string{
	1:  "CONNECT",
	2:  "CONNACK",
	3:  "PUBLISH",
	4:  "PUBACK",
	5:  "PUBREC",
	6:  "PUBREL",
	7:  "PUBCOMP",
	8:  "SUBSCRIBE",
	9:  "SUBACK",
	10: "UNSUBSCRIBE",
	11: "UNSUBACK",
	12: "PINGREQ",
	13: "PINGRESP",
	14: "DISCONNECT",
}