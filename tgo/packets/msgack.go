package packets

import "fmt"

type MsgackPacket struct {
	FixedHeader
	MessageID uint64
}

func NewMsgackPacketWithHeader(fh FixedHeader) *MsgackPacket  {
	m := &MsgackPacket{}
	m.FixedHeader = fh
	return  m
}
func NewMsgackPacket(messageID uint64) *MsgackPacket {
	m := &MsgackPacket{}
	m.PacketType = Msgack
	m.MessageID = messageID
	return m
}

func (m *MsgackPacket) GetFixedHeader() FixedHeader  {

	return m.FixedHeader
}

func (m *MsgackPacket) String() string {
	str := fmt.Sprintf("%s", m.FixedHeader)
	str += " "
	str += fmt.Sprintf("MessageID: %d", m.MessageID)
	return str
}