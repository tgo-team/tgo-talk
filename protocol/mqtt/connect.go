package mqtt

import (
	"bytes"
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"io"
)

func (m *MQTTCodec) decodeConnect(fh *packets.FixedHeader, reader io.Reader) (*packets.ConnectPacket, error) {
	c := packets.NewConnectPacket(*fh)
	var _ = decodeString(reader)
	var _ = decodeByte(reader)
	options := decodeByte(reader)
	_ = 1 & options        // ReservedBit
	_ = 1&(options>>1) > 0 // CleanSession
	_ = 1&(options>>2) > 0 // WillFlag
	_ = 3 & (options >> 3) // WillQos
	_ = 1&(options>>5) > 0 //WillRetain
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0
	c.Keepalive = decodeUint16(reader)
	c.ClientIdentifier = decodeUint64(reader)
	if c.UsernameFlag {
		c.Username = decodeString(reader)
	}
	if c.PasswordFlag {
		c.Password = decodeBytes(reader)
	}
	return c, nil
}

func (m *MQTTCodec) encodeConnect(packet packets.Packet) ([]byte, error) {
	c := packet.(*packets.ConnectPacket)
	var body bytes.Buffer
	body.Write(encodeString("MQTT"))
	body.WriteByte(0x04)
	body.WriteByte(boolToByte(false)<<1 | boolToByte(false)<<2 | 0<<3 | boolToByte(false)<<5 | boolToByte(c.PasswordFlag)<<6 | boolToByte(c.UsernameFlag)<<7)
	body.Write(encodeUint16(c.Keepalive))
	body.Write(encodeUint64(c.ClientIdentifier))
	if c.UsernameFlag {
		body.Write(encodeString(c.Username))
	}
	if c.PasswordFlag {
		body.Write(encodeBytes(c.Password))
	}
	c.RemainingLength = body.Len()
	return body.Bytes(), nil
}
