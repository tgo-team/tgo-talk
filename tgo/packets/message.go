package packets

import "fmt"

type MessagePacket struct {
	FixedHeader
	GIDFlag bool  // 是否存在GID
	UID uint64 // 接受用户的ID
	GID uint64 // 群组的ID
	MessageID uint64
	Payload   []byte
}

func NewMessagePacket(fh FixedHeader) *MessagePacket  {
	p := &MessagePacket{}
	p.FixedHeader = fh
	return  p
}


func (p *MessagePacket) GetFixedHeader() FixedHeader  {

	return p.FixedHeader
}

func (p *MessagePacket) String() string {
	str := fmt.Sprintf("%s", p.FixedHeader)
	str += " "
	str += fmt.Sprintf("UID: %d GID: %d MessageID: %d", p.UID,p.GID, p.MessageID)
	str += " "
	str += fmt.Sprintf("payload: %s", string(p.Payload))
	return str
}