package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeCMD(fh *packets.FixedHeader,reader io.Reader) (*packets.CmdPacket, error) {
	c :=packets.NewCmdPacketWithHeader(*fh)
	c.CMD = packets.DecodeString(reader)
	var payloadLength = c.RemainingLength - (len(c.CMD) + 2) // payloadLength = 剩余长度 - CMD长度
	if payloadLength < 0 {
		return nil,fmt.Errorf("Error upacking cmd, payload length < 0")
	}
	c.Payload = make([]byte, payloadLength)
	_, err := reader.Read(c.Payload)
	return c,err
}

func (m *MQTTCodec) encodeCMD(packet packets.Packet) ([]byte, error) {
	c := packet.(*packets.CmdPacket)
	var body bytes.Buffer
	body.Write(packets.EncodeString(c.CMD))
	body.Write(c.Payload)
	c.RemainingLength = body.Len()
	return body.Bytes(),nil
}
