package tcp

import (
	"errors"
	_ "github.com/tgo-team/tgo-core/protocol/mqtt"
	_ "github.com/tgo-team/tgo-core/storage/memory"
	"github.com/tgo-team/tgo-core/test"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"net"
	"testing"
	"time"
)

func TestTCPServer_StartAndStop(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	tg := startTGO(opts)
	time.Sleep(50 * time.Millisecond)
	err := tg.Stop()
	test.Nil(t, err)
	time.Sleep(50 * time.Millisecond)
}

func TestTCPServer_ReceivePacketChan(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	tg := startTGO(opts)

	clientConn, err := MustConnectServer(tg.Server.(*TCPServer).RealTCPAddr())
	test.Nil(t, err)

	connectData, err := tg.GetOpts().Pro.EncodePacket(&packets.ConnectPacket{FixedHeader: packets.FixedHeader{PacketType: packets.Connect}, ClientIdentifier: 1, PasswordFlag: true, Password: []byte("123456")})
	test.Nil(t, err)
	_, err = clientConn.Write(connectData)
	test.Nil(t, err)

	time.Sleep(50 * time.Millisecond)

	tg.Stop()
}

func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func startTGO(opts *tgo.Options) *tgo.TGO {
	opts.TCPAddress = "127.0.0.1:0"
	opts.HTTPAddress = "127.0.0.1:0"
	opts.HTTPSAddress = "127.0.0.1:0"
	tg := tgo.New(opts)
	err := tg.Start()
	if err != nil {
		panic(err)
	}
	return tg
}

type authTest struct {
}

func (a *authTest) Auth(clientID uint64, password string) error {
	if clientID == 1 && password == "123456" {
		return nil
	}
	return errors.New("认证错误！")
}
