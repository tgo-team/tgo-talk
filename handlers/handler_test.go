package handlers

import (
	_ "github.com/tgo-team/tgo-talk/log"
	_ "github.com/tgo-team/tgo-talk/protocol/mqtt"
	"github.com/tgo-team/tgo-talk/server/tcp"
	_ "github.com/tgo-team/tgo-talk/storage/memory"
	"github.com/tgo-team/tgo-talk/test"
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"net"
	"testing"
	"time"
)

func TestHandle(t *testing.T) {
	tg := startTGO(t)
	packetChan := make(chan packets.Packet, 0)
	tg.Use(func(context *tgo.MContext) {
		packetChan <- context.Packet()
	})

	conn, err := MustConnectServer(tg.Server.(*tcp.Server).RealTCPAddr())
	test.Nil(t, err)
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientIdentifier: 1, PasswordFlag: true, Password: []byte("123456")}
	WritePacket(t, conn, connPacket, tg)

	packet := <-packetChan
	test.Equal(t, connPacket.String(), packet.String())
}

func TestHandleAuth(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleAuth)

	conn, err := MustConnectServer(tg.Server.(*tcp.Server).RealTCPAddr())
	test.Nil(t, err)
	sendAuthPacket(t, conn, tg)

}

func TestHandleHeartbeat(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleAuth)
	tg.Use(HandleHeartbeat)

	conn, err := MustConnectServer(tg.Server.(*tcp.Server).RealTCPAddr())
	test.Nil(t, err)

	sendAuthPacket(t, conn, tg)

	// 心跳
	pingPacket := &packets.PingreqPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Pingreq}}
	WritePacket(t, conn, pingPacket, tg)

	pongPacket, ok := ReadPacket(t, conn, tg).(*packets.PingrespPacket)
	test.Equal(t, true, ok)
	test.Equal(t, packets.Pingresp, pongPacket.GetFixedHeader().PacketType)
}

func TestHandleRevMsg(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleAuth)
	tg.Use(HandleHeartbeat)
	tg.Use(HandleRevMsg)

	conn, err := MustConnectServer(tg.Server.(*tcp.Server).RealTCPAddr())
	test.Nil(t, err)
	sendAuthPacket(t, conn, tg)

	sendMsgPacket(t,conn,tg,packets.NewMessagePacket(100,2,[]byte("hello")))

	time.Sleep(time.Millisecond*50)

}

func sendMsgPacket(t *testing.T, conn net.Conn, tg *tgo.TGO, msgPacket *packets.MessagePacket) {
	WritePacket(t, conn, msgPacket, tg)
}

func sendAuthPacket(t *testing.T, conn net.Conn, tg *tgo.TGO) {
	// 认证
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientIdentifier: 1, PasswordFlag: true, Password: []byte("123456")}
	WritePacket(t, conn, connPacket, tg)
	connackPacket, ok := ReadPacket(t, conn, tg).(*packets.ConnackPacket)

	test.Equal(t, true, ok)
	test.Equal(t, packets.Connack, connackPacket.GetFixedHeader().PacketType)
	test.Equal(t, packets.ConnReturnCodeSuccess, connackPacket.ReturnCode)
}

func WritePacket(t *testing.T, conn net.Conn, packet packets.Packet, tg *tgo.TGO) {
	pingData, err := tg.GetOpts().Pro.EncodePacket(packet)
	test.Nil(t, err)
	_, err = conn.Write(pingData)
}

func ReadPacket(t *testing.T, conn net.Conn, tg *tgo.TGO) packets.Packet {
	packetObj, err := tg.GetOpts().Pro.DecodePacket(conn)
	test.Nil(t, err)
	return packetObj
}

func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func startTGO(t *testing.T) *tgo.TGO {
	opts := tgo.NewOptions()
	opts.TCPAddress = "0.0.0.0:0"
	opts.Log = test.NewLog(t)
	tg := tgo.New(opts)
	err := tg.Start()
	test.Nil(t, err)
	return tg
}
