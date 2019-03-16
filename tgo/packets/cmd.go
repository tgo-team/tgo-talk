package packets

import "fmt"

type CMDPacket struct {
	FixedHeader
	CMD       uint16
	Payload   []byte // 消息内容
}

func NewCMDPacketWithHeader(fh FixedHeader) *CMDPacket {
	c := &CMDPacket{}
	c.FixedHeader = fh
	return c
}

func NewCMDPacket(cmd uint16, payload []byte) *CMDPacket {

	return &CMDPacket{ CMD: cmd, Payload: payload,FixedHeader:FixedHeader{PacketType:CMD}}
}

func (c *CMDPacket) GetFixedHeader() FixedHeader {

	return c.FixedHeader
}

func (c *CMDPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("CMD: %d Payload:  %s", c.CMD, string(c.Payload))
	return str
}
