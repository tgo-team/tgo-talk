package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeMessage(fh *packets.FixedHeader,reader io.Reader) ( *packets.MessagePacket, error) {
	msg := packets.NewMessagePacketHeader(*fh)
	var payloadLength = msg.RemainingLength
	msg.From = packets.DecodeUint64(reader)
	if msg.From == 0 {
		statefulConn,ok := reader.(tgo.StatefulConn)
		if ok {
			msg.From = statefulConn.GetID()
		}
	}
	msg.ChannelID = packets.DecodeUint64(reader)
	payloadLength -= 8 + 8 // 减去 ChannelID的长度 + From的长度
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
	if msg.Qos > 0 {
		body.Write(packets.EncodeUint64(msg.MessageID))
	}
	body.Write(msg.Payload)
	msg.RemainingLength = body.Len()
	return body.Bytes(),nil
}