package mqtt

import (
	"bytes"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeConnect(fh *packets.FixedHeader, reader io.Reader) (*packets.ConnectPacket, error) {
	c := packets.NewConnectPacketWithHeader(*fh)
	var _ = packets.DecodeString(reader)
	var _ = packets.DecodeByte(reader)
	options := packets.DecodeByte(reader)
	_ = 1 & options        // ReservedBit
	_ = 1&(options>>1) > 0 // CleanSession
	_ = 1&(options>>2) > 0 // WillFlag
	_ = 3 & (options >> 3) // WillQos
	_ = 1&(options>>5) > 0 //WillRetain
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0
	c.Keepalive = packets.DecodeUint16(reader)
	c.ClientID = packets.DecodeUint64(reader)
	if c.UsernameFlag {
		c.Username = packets.DecodeString(reader)
	}
	if c.PasswordFlag {
		c.Password = string(packets.DecodeBytes(reader))
	}
	return c, nil
}

func (m *MQTTCodec) encodeConnect(packet packets.Packet) ([]byte, error) {
	c := packet.(*packets.ConnectPacket)
	var body bytes.Buffer
	body.Write(packets.EncodeString("MQTT"))
	body.WriteByte(0x04)
	body.WriteByte(packets.BoolToByte(false)<<1 | packets.BoolToByte(false)<<2 | 0<<3 | packets.BoolToByte(false)<<5 | packets.BoolToByte(c.PasswordFlag)<<6 | packets.BoolToByte(c.UsernameFlag)<<7)
	body.Write(packets.EncodeUint16(c.Keepalive))
	body.Write(packets.EncodeUint64(c.ClientID))
	if c.UsernameFlag {
		body.Write(packets.EncodeString(c.Username))
	}
	if c.PasswordFlag {
		body.Write(packets.EncodeString(c.Password))
	}
	c.RemainingLength = body.Len()
	return body.Bytes(), nil
}
