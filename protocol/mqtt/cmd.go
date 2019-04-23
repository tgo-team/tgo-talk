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
	options := packets.DecodeByte(reader)
	println( options>>7)
	c.TokenFlag = 1 &(options>>7) > 0 // 是否有token
	if c.TokenFlag {
		c.Token = packets.DecodeString(reader)
	}
	var payloadLength = c.RemainingLength - (len(c.CMD) + 2)  - 1// payloadLength = 剩余长度 - CMD长度 - 减去option的1byte
	if c.TokenFlag {
		payloadLength = payloadLength - len(c.Token) - 2 // 减去token字符串长度和字符串占位的2位
	}
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
	body.WriteByte(packets.BoolToByte(c.TokenFlag)<<7)
	if c.TokenFlag {
		body.Write(packets.EncodeString(c.Token))
	}
	body.Write(c.Payload)
	c.RemainingLength = body.Len()
	return body.Bytes(),nil
}
