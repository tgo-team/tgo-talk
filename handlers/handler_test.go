package handlers

import (
	_ "github.com/tgo-team/tgo-chat/log"
	_ "github.com/tgo-team/tgo-chat/protocol"
	"github.com/tgo-team/tgo-chat/server"
	_ "github.com/tgo-team/tgo-chat/server"
	_ "github.com/tgo-team/tgo-chat/storage/memory"
	"github.com/tgo-team/tgo-chat/test"
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	tg := tgo.New(opts)
	hasValue := 0
	tg.Use(func(context *tgo.MContext) {
		hasValue = 1
	})
	err := tg.Start()
	test.Nil(t, err)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)

	_, err = cn.Write([]byte{0x06})
	test.Nil(t, err)

	time.Sleep(time.Millisecond * 50)

	test.Equal(t, 1, hasValue)
}

func TestHandleHeartbeat(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	tg := tgo.New(opts)
	tg.Use(HandleHeartbeat)
	err := tg.Start()
	test.Nil(t, err)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)

	_, err = cn.Write([]byte{byte(tgo.MsgTypePing)})
	test.Nil(t, err)
	time.Sleep(time.Millisecond * 50)


	resultBytes := make([]byte,1)
	_,err = cn.Read(resultBytes)
	test.Nil(t,err)

	test.Equal(t,tgo.MsgTypePong,tgo.MsgType(resultBytes[0]))
}

func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
