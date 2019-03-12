package tcp

import (
	_ "github.com/tgo-team/tgo-talk/protocol/mqtt"
	"github.com/tgo-team/tgo-talk/test"
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"net"
	"testing"
)

func TestNewConn(t *testing.T) {
	c, s := net.Pipe()
	packConnChan := make(chan *tgo.PacketConn, 0)
	connExitChan := make(chan tgo.Conn, 0)
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	NewConn(1000100010001000100, s, packConnChan, connExitChan, opts)
	connPacket := &packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientIdentifier: 1, PasswordFlag: true, Password: []byte("123456")}
	connectData, err := opts.Pro.EncodePacket(connPacket)
	test.Nil(t, err)

	_, err = c.Write(connectData)
	test.Nil(t, err)

	packetConn := <-packConnChan
	test.Equal(t, connPacket.String(), packetConn.Packet.String())
}
