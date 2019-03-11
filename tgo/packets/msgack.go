package packets

import "fmt"

type MsgackPacket struct {
	FixedHeader
	MessageID uint64
}

func NewMsgackPacket(fh FixedHeader) *MsgackPacket  {
	m := &MsgackPacket{}
	m.FixedHeader = fh
	return  m
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