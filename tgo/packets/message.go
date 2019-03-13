package packets

import "fmt"

type MessagePacket struct {
	FixedHeader
	ChannelID uint64 // 管道ID
	MessageID uint64 // 消息唯一编号
	Payload   []byte // 消息内容
}

func NewMessagePacket(messageID uint64, channelID uint64, payload []byte) *MessagePacket {
	return &MessagePacket{
		ChannelID:   channelID,
		MessageID:   messageID,
		Payload:     payload,
		FixedHeader: FixedHeader{PacketType: Message, Qos: 1},
	}
}

func NewMessagePacketHeader(fh FixedHeader) *MessagePacket {
	p := &MessagePacket{}
	p.FixedHeader = fh
	return p
}

func (p *MessagePacket) GetFixedHeader() FixedHeader {

	return p.FixedHeader
}

func (p *MessagePacket) String() string {
	str := fmt.Sprintf("%s", p.FixedHeader)
	str += " "
	str += fmt.Sprintf("ChannelID: %d MessageID: %d", p.ChannelID, p.MessageID)
	str += " "
	str += fmt.Sprintf("payload: %s", string(p.Payload))
	return str
}
