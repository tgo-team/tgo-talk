package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeCMD(fh *packets.FixedHeader,reader io.Reader) (*packets.CMDPacket, error) {
	c :=packets.NewCMDPacketWithHeader(*fh)
	c.CMD = packets.DecodeUint16(reader)
	var payloadLength = c.RemainingLength - 2 - 8 // payloadLength = 剩余长度 - CMD长度 - MessageID长度
	if payloadLength < 0 {
		return nil,fmt.Errorf("Error upacking cmd, payload length < 0")
	}
	c.Payload = make([]byte, payloadLength)
	_, err := reader.Read(c.Payload)
	return c,err
}

func (m *MQTTCodec) encodeCMD(packet packets.Packet) ([]byte, error) {
	c := packet.(*packets.CMDPacket)
	var body bytes.Buffer
	body.Write(packets.EncodeUint16(c.CMD))
	body.Write(c.Payload)
	c.RemainingLength = body.Len()
	return body.Bytes(),nil
}
