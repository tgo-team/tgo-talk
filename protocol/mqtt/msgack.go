package mqtt

import (
	"bytes"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeMsgack(fh *packets.FixedHeader,reader io.Reader) ( *packets.MsgackPacket, error) {
	msg := packets.NewMsgackPacketWithHeader(*fh)
	messageIDCount := msg.RemainingLength/8
	messageIDs := make([]uint64,0)
	for i:=0;i<messageIDCount;i++ {
		messageIDs = append(messageIDs,packets.DecodeUint64(reader))
	}
	msg.MessageIDs = messageIDs
	return msg,nil
}

func (m *MQTTCodec) encodeMsgack(packet packets.Packet) ([]byte, error) {
	msg := packet.(*packets.MsgackPacket)
	var body bytes.Buffer
	for _,messageID :=range msg.MessageIDs {
		body.Write(packets.EncodeUint64(messageID))
	}
	msg.RemainingLength = body.Len()
	return body.Bytes(),nil
}
