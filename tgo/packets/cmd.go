package packets

import "fmt"

type CMDPacket struct {
	FixedHeader
	CMD       uint16
	MessageID uint64 // 消息唯一编号
	Payload   []byte // 消息内容
}

func NewCMDPacketWithHeader(fh FixedHeader) *CMDPacket {
	c := &CMDPacket{}
	c.FixedHeader = fh
	return c
}

func NewCMDPacket(messageID uint64, cmd uint16, payload []byte) *CMDPacket {

	return &CMDPacket{MessageID: messageID, CMD: cmd, Payload: payload,FixedHeader:FixedHeader{PacketType:CMD}}
}

func (c *CMDPacket) GetFixedHeader() FixedHeader {

	return c.FixedHeader
}

func (c *CMDPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d CMD: %d Payload:  %s", c.MessageID, c.CMD, string(c.Payload))
	return str
}
