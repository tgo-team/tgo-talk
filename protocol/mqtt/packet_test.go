package mqtt

import (
	"bytes"
	"github.com/tgo-team/tgo-talk/test"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"testing"
)

func TestMQTTCodec_Decode_ConnectPacket(t *testing.T)  {
	connectPacketBytes := bytes.NewBuffer([]byte{16, 52, 0, 4, 77, 81, 84, 84, 4, 192, 0, 0, 0, 0, 0, 8, 116, 101, 115, 116, 117, 115, 101, 114, 0, 8, 116, 101, 115, 116, 112, 97, 115, 115})
	codec := &MQTTCodec{}
	packet,err := codec.DecodePacket(connectPacketBytes)
	cp := packet.(*packets.ConnectPacket)
	test.Nil(t,err)
	if cp.UsernameFlag != true {
		t.Errorf("Connect Packet UsernameFlag is %t, should be %t", cp.UsernameFlag, true)
	}
	if cp.Username != "testuser" {
		t.Errorf("Connect Packet Username is %s, should be %s", cp.Username, "testuser")
	}
	if cp.PasswordFlag != true {
		t.Errorf("Connect Packet PasswordFlag is %t, should be %t", cp.PasswordFlag, true)
	}
	if string(cp.Password) != "testpass" {
		t.Errorf("Connect Packet Password is %s, should be %s", string(cp.Password), "testpass")
	}
}

func TestMQTTCodec_Encode_ConnectPacket(t *testing.T)  {

	connectPacket := &packets.ConnectPacket{}
	connectPacket.UsernameFlag = true
	connectPacket.PasswordFlag = true
	connectPacket.Username = "testuser"
	connectPacket.Password = "testpass"
	codec := &MQTTCodec{}
	_,err := codec.EncodePacket(connectPacket)
	test.Nil(t,err)
}

func TestMQTTCodec_DecodeAndEncode_ConnectPacket(t *testing.T)  {
	data := []byte{16, 32, 0, 4, 77, 81, 84, 84, 4, 192, 0, 0, 0, 0, 0, 8, 116, 101, 115, 116, 117, 115, 101, 114, 0, 8, 116, 101, 115, 116, 112, 97, 115, 115}
	connectPacketBuff := bytes.NewBuffer(data)
	codec := &MQTTCodec{}
	packet,err := codec.DecodePacket(connectPacketBuff)
	cp := packet.(*packets.ConnectPacket)
	test.Nil(t,err)

	b,err := codec.EncodePacket(cp)

	test.Equal(t,data,b)
}

func TestMQTTCodec_DecodeAndEncodePackets(t *testing.T) {
	ps := []packets.Packet{
		&packets.ConnectPacket{FixedHeader:packets.FixedHeader{PacketType:packets.Connect},Keepalive:10},
		&packets.ConnackPacket{FixedHeader:packets.FixedHeader{PacketType:packets.Connack},ReturnCode:1},
		&packets.PingreqPacket{FixedHeader:packets.FixedHeader{PacketType:packets.Pingreq}},
		&packets.PingrespPacket{FixedHeader:packets.FixedHeader{PacketType:packets.Pingresp}},
		&packets.MessagePacket{FixedHeader:packets.FixedHeader{Qos:1,PacketType:packets.Message},ChannelID:123,MessageID:234,Payload:[]byte("hello")},
		&packets.MsgackPacket{FixedHeader:packets.FixedHeader{PacketType:packets.Msgack},MessageID:234},
	}
	codec := &MQTTCodec{}
	for _,read :=range ps {
		b,err := codec.EncodePacket(read)
		test.Nil(t,err)
		packet,err := codec.DecodePacket(bytes.NewBuffer(b))
		test.Nil(t,err)
		println(packet.String())
		println(read.String())
		test.Equal(t,packet.String(),read.String())
	}
}