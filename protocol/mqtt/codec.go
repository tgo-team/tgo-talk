package mqtt

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"io"
)

func init() {
	tgo.RegistryProtocol("mqtt-im", func() tgo.Protocol {
		return NewMQTTCodec()
	})
}

type MQTTCodec struct {
}

func NewMQTTCodec() *MQTTCodec {
	return &MQTTCodec{}
}

func (m *MQTTCodec) DecodePacket(reader tgo.Conn) (packets.Packet, error) {
	fh, err := m.decodeFixedHeader(reader)
	if err != nil {
		return nil, err
	}

	if fh.PacketType == packets.Connect {
		return m.decodeConnect(fh, reader)
	}
	if fh.PacketType == packets.Connack {
		return m.decodeConnack(fh, reader)
	}
	if fh.PacketType == packets.Pingreq {
		return &packets.PingreqPacket{FixedHeader: *fh}, nil
	}
	if fh.PacketType == packets.Pingresp {
		return &packets.PingrespPacket{FixedHeader: *fh}, nil
	}
	if fh.PacketType == packets.Message {
		return  m.decodeMessage(fh, reader)
	}
	if fh.PacketType == packets.Msgack {
		return m.decodeMsgack(fh, reader)
	}
	if fh.PacketType == packets.Cmd {
		return m.decodeCMD(fh, reader)
	}
	return nil, fmt.Errorf("不支持的包类型[%d]",fh.PacketType)
}

func (m *MQTTCodec) EncodePacket(packet packets.Packet) ([]byte, error) {

	var packetType = packet.GetFixedHeader().PacketType

	var packetBuffer bytes.Buffer
	var remainingBytes []byte
	var err error

	if packetType == packets.Connect {
		remainingBytes, err = m.encodeConnect(packet)
		if err != nil {
			return nil, err
		}
	}
	if packetType == packets.Connack {
		remainingBytes, err = m.encodeConnack(packet)
		if err != nil {
			return nil, err
		}
	}
	if packetType == packets.Message {
		remainingBytes, err = m.encodeMessage(packet)
		if err != nil {
			return nil, err
		}
	}
	if packetType == packets.Msgack {
		remainingBytes, err = m.encodeMsgack(packet)
		if err != nil {
			return nil, err
		}
	}
	if packetType == packets.Cmd {
		remainingBytes, err = m.encodeCMD(packet)
		if err != nil {
			return nil, err
		}
	}

	header, err := m.encodeFixedHeader(packet)
	if err != nil {
		return nil, err
	}
	if packetType == packets.Pingreq || packetType == packets.Pingresp {
		return header, nil
	}

	if remainingBytes ==nil || len(remainingBytes)<=0 {
		return nil,fmt.Errorf("不支持的包类型[%d]",packetType)
	}

	packetBuffer.Write(header)
	packetBuffer.Write(remainingBytes)

	return packetBuffer.Bytes(), nil
}

func (m *MQTTCodec) decodeFixedHeader(reader tgo.Conn) (*packets.FixedHeader, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(reader, b)
	if err != nil {
		return nil, err
	}
	typeAndFlags := b[0]
	fh := &packets.FixedHeader{}
	fh.PacketType = packets.PacketType(typeAndFlags >> 4)
	fh.Dup = (typeAndFlags>>3)&0x01 > 0
	fh.Qos = (typeAndFlags >> 1) & 0x03
	fh.Retain = typeAndFlags&0x01 > 0
	fh.RemainingLength = decodeLength(reader)

	statefulCon,ok := reader.(tgo.StatefulConn)
	if ok {
		fh.From = statefulCon.GetID()
	}

	return fh, nil
}

func (m *MQTTCodec) encodeFixedHeader(packet packets.Packet) ([]byte, error) {
	fh := packet.GetFixedHeader()
	var header bytes.Buffer
	header.WriteByte(byte(packet.GetFixedHeader().PacketType)<<4 | packets.BoolToByte(fh.Dup)<<3 | fh.Qos<<1 | packets.BoolToByte(fh.Retain))
	header.Write(encodeLength(fh.RemainingLength))
	return header.Bytes(), nil
}
