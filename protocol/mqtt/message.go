package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"io"
)

// | Option(1 byte ) |
//  | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |
//  | gid flag |
func (m *MQTTCodec) decodeMessage(fh *packets.FixedHeader,reader io.Reader) ( *packets.MessagePacket, error) {
	msg := packets.NewMessagePacket(*fh)
	var payloadLength = msg.RemainingLength
	options := decodeByte(reader)
	msg.GIDFlag = (options>>7) >0
	msg.UID = decodeUint64(reader)
	payloadLength -= 8 + 1// 减去 uid的长度 + options的长度
	if msg.GIDFlag {
		msg.GID = decodeUint64(reader)
		payloadLength -= 8 // 减去gid长度
	}
	if msg.Qos > 0 {
		msg.MessageID = decodeUint64(reader)
		payloadLength -=  8 // 减去messageID长度
	}
	if payloadLength < 0 {
		return nil,fmt.Errorf("Error upacking publish, payload length < 0")
	}
	msg.Payload = make([]byte, payloadLength)
	_, err := reader.Read(msg.Payload)
	return msg,err
}

func (m *MQTTCodec) encodeMessage(packet packets.Packet) ([]byte, error) {
	msg := packet.(*packets.MessagePacket)
	var body bytes.Buffer
	body.WriteByte(boolToByte(msg.GIDFlag)<<7)
	body.Write(encodeUint64(msg.UID))
	if msg.GIDFlag {
		body.Write(encodeUint64(msg.GID))
	}
	if msg.Qos > 0 {
		body.Write(encodeUint64(msg.MessageID))
	}
	body.Write(msg.Payload)
	msg.RemainingLength = body.Len()
	return body.Bytes(),nil
}