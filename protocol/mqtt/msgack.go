package mqtt

import (
	"bytes"
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeMsgack(fh *packets.FixedHeader,reader io.Reader) ( *packets.MsgackPacket, error) {
	msg := packets.NewMsgackPacket(*fh)
	msg.MessageID = decodeUint64(reader)
	return msg,nil
}

func (m *MQTTCodec) encodeMsgack(packet packets.Packet) ([]byte, error) {
	msg := packet.(*packets.MsgackPacket)
	var body bytes.Buffer
	body.Write(encodeUint64(msg.MessageID))
	msg.RemainingLength = body.Len()
	return body.Bytes(),nil
}
