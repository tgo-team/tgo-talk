package packets

import "fmt"

type PingrespPacket struct {
	FixedHeader
}

func NewPingrespPacketWithHeader(fh FixedHeader) *PingrespPacket  {
	pr := &PingrespPacket{}
	pr.FixedHeader = fh
	return pr
}

func NewPingrespPacket() *PingrespPacket {
	pr := &PingrespPacket{}
	pr.PacketType = Pingresp
	return pr
}

func (pr *PingrespPacket) GetFixedHeader() FixedHeader  {

	return pr.FixedHeader
}

func (pr *PingrespPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}
