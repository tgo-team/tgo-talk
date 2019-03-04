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

func TestHandleAuth(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	tg := tgo.New(opts)
	tg.Use(HandleAuth)
	tg.Use(HandleHeartbeat)
	err := tg.Start()
	test.Nil(t, err)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)
	msgBytes := []byte{byte(tgo.MsgTypeAuth),13,0x00,0x00,0x00,0x03,0x03,0x03,0x03,0x03,0x03,0x03,0x03,0,0}
	msgBytes = append(msgBytes,[]byte("pwd")...)
	_, err = cn.Write(msgBytes)
	test.Nil(t, err)
	time.Sleep(time.Millisecond * 50)


	msg,err := tg.GetOpts().Pro.Decode(cn)
	test.Nil(t,err)

	test.Equal(t,tgo.MsgTypeAuthACK,msg.MsgType)
	test.Equal(t,tgo.MsgStatusAuthOk,int8(msg.VariableHeader[0]))
}

func authData(uid int64,token string) []byte {

}

func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
