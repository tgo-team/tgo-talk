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
	connectData, err := tg.GetOpts().Pro.EncodePacket(connPacket)

	_, err = conn.Write(connectData)
	test.Nil(t, err)

	packet := <-packetChan
	test.Equal(t, connPacket.String(), packet.String())
}

func TestHandleAuth(t *testing.T) {
	tg := startTGO(t)
	tg.Use(HandleAuth)

	conn, err := MustConnectServer(tg.Server.(*tcp.Server).RealTCPAddr())
	test.Nil(t, err)
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientIdentifier: 1, PasswordFlag: true, Password: []byte("123456")}
	connectData, err := tg.GetOpts().Pro.EncodePacket(connPacket)

	_, err = conn.Write(connectData)
	test.Nil(t, err)

	connackPacketObj, err := tg.GetOpts().Pro.DecodePacket(conn)
	test.Nil(t, err)

	time.Sleep(time.Millisecond*50)

	connackPacket, ok := connackPacketObj.(*packets.ConnackPacket)
	test.Equal(t, true, ok)
	test.Equal(t, packets.Connack, connackPacket.GetFixedHeader().PacketType)
	test.Equal(t, packets.ConnReturnCodeSuccess, connackPacket.ReturnCode)

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
