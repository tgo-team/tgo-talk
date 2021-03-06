package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeMessage(fh *packets.FixedHeader,reader io.Reader) ( *packets.MessagePacket, error) {
	msg := packets.NewMessagePacketHeader(*fh)
	var payloadLength = msg.RemainingLength
	from := packets.DecodeUint64(reader)
	if msg.From == 0 { // 如果不存在from 则使用协议里的 否则使用header里的from （因为有状态连接不需要从协议里获取from）
		msg.From = from
	}
	msg.ChannelID = packets.DecodeUint64(reader)
	msg.Timestamp = int64(packets.DecodeUint32(reader))
	payloadLength -= 8 + 8 + 4 // 减去 From的长度 + ChannelID的长度 + Timestamp的长度
	if msg.Qos > 0 {
		msg.MessageID = packets.DecodeUint64(reader)
		payloadLength -=  8 // 减去messageID长度
	}
	if payloadLength < 0 {
		return nil,fmt.Errorf("Error upacking message, payload length < 0")
	}
	msg.Payload = make([]byte, payloadLength)
	_, err := reader.Read(msg.Payload)
	return msg,err
}

func (m *MQTTCodec) encodeMessage(packet packets.Packet) ([]byte, error) {
	msg := packet.(*packets.MessagePacket)
	var body bytes.Buffer
	body.Write(packets.EncodeUint64(msg.From))
	body.Write(packets.EncodeUint64(msg.ChannelID))
	body.Write(packets.EncodeUint32(uint32(msg.Timestamp)))
	if msg.Qos > 0 {
		body.Write(packets.EncodeUint64(msg.MessageID))
	}
	body.Write(msg.Payload)
	msg.RemainingLength = body.Len()
	return body.Bytes(),nil
}