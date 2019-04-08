package handlers

import (
	"fmt"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
	_ "github.com/tgo-team/tgo-talk/log"
	_ "github.com/tgo-team/tgo-talk/protocol/mqtt"
	"github.com/tgo-team/tgo-talk/server/tcp"
	_ "github.com/tgo-team/tgo-talk/storage/memory"
	"github.com/tgo-team/tgo-talk/test"
	"net"
	"os"
	"testing"
	"time"
)

var testClientID uint64 = 1
var testPassword = "123456"
var testChannelID uint64 = 1

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestHandle(t *testing.T) {
	tg := startTGO(t)
	packetChan := make(chan packets.Packet, 0)
	tg.Use(func(context *tgo.MContext) {
		packetChan <- context.Packet()
	})
	var tcpServer *tcp.Server
	for _, server := range tg.Servers {
		s, ok := server.(*tcp.Server)
		if ok {
			tcpServer = s
		}
	}
	conn, err := MustConnectServer(tcpServer.RealTCPAddr())
	test.Nil(t, err)
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientID: 1, PasswordFlag: true, Password: "123456"}
	WritePacket(t, conn, connPacket, tg)

	packet := <-packetChan
	test.Equal(t, connPacket.String(), packet.String())
}

func TestHandleConnPacket(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleConnPacket)

	var tcpServer *tcp.Server
	for _, server := range tg.Servers {
		s, ok := server.(*tcp.Server)
		if ok {
			tcpServer = s
		}
	}

	conn, err := MustConnectServer(tcpServer.RealTCPAddr())
	test.Nil(t, err)
	sendAuthPacket(t, conn, tg)

}

func TestHandlePingPacket(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleConnPacket)
	tg.Use(HandlePingPacket)
	var tcpServer *tcp.Server
	for _, server := range tg.Servers {
		s, ok := server.(*tcp.Server)
		if ok {
			tcpServer = s
		}
	}

	conn, err := MustConnectServer(tcpServer.RealTCPAddr())
	test.Nil(t, err)

	sendAuthPacket(t, conn, tg)

	// 心跳
	pingPacket := &packets.PingreqPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Pingreq}}
	WritePacket(t, conn, pingPacket, tg)

	pongPacket, ok := ReadPacket(t, conn, tg).(*packets.PingrespPacket)
	test.Equal(t, true, ok)
	test.Equal(t, packets.Pingresp, pongPacket.GetFixedHeader().PacketType)
}

func TestHandleMessagePacket(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleConnPacket)
	tg.Use(HandlePingPacket)
	tg.Match(fmt.Sprintf("type:%d", packets.Message), HandleMessagePacket)

	var tcpServer *tcp.Server
	for _, server := range tg.Servers {
		s, ok := server.(*tcp.Server)
		if ok {
			tcpServer = s
		}
	}
	conn, err := MustConnectServer(tcpServer.RealTCPAddr())
	test.Nil(t, err)
	sendAuthPacket(t, conn, tg)

	sendMsgPacket(t, conn, tg, packets.NewMessagePacket(100, testChannelID, []byte("hello")))

	msgackPacket, ok := ReadPacket(t, conn, tg).(*packets.MsgackPacket)
	test.Equal(t, true, ok)
	test.Equal(t, packets.Msgack, msgackPacket.PacketType)
	test.Equal(t, uint64(100), msgackPacket.MessageID)

}

func TestHandleCmdPacket(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleConnPacket)
	tg.Use(HandlePingPacket)
	tg.Match(fmt.Sprintf("type:%d", packets.CMD), HandleCmdPacket)
	var tcpServer *tcp.Server
	for _, server := range tg.Servers {
		s, ok := server.(*tcp.Server)
		if ok {
			tcpServer = s
		}
	}
	conn, err := MustConnectServer(tcpServer.RealTCPAddr())
	test.Nil(t, err)
	sendCmdPacket(t, conn, tg, packets.NewCmdPacket(100, []byte("admin")))
	time.Sleep(time.Millisecond * 50)
}

func sendMsgPacket(t *testing.T, conn net.Conn, tg *tgo.TGO, msgPacket *packets.MessagePacket) {
	WritePacket(t, conn, msgPacket, tg)
}

func sendCmdPacket(t *testing.T, conn net.Conn, tg *tgo.TGO, CmdPacket *packets.CmdPacket) {
	WritePacket(t, conn, CmdPacket, tg)
}

func sendAuthPacket(t *testing.T, conn net.Conn, tg *tgo.TGO) {
	// 认证
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientID: testClientID, PasswordFlag: true, Password: testPassword}
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
	packetObj, err := tg.GetOpts().Pro.DecodePacket(tcp.NewConn(conn, nil, nil))
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
	tg.Storage.AddClient(tgo.NewClient(testClientID, testPassword))
	tg.Storage.AddChannel(tgo.NewChannel(testChannelID, 1, &tgo.Context{TGO: tg}))
	test.Nil(t, err)
	return tg
}
