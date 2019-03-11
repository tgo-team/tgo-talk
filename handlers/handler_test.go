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
	tg := startTGO(t)
	hasValue := 0
	tg.Use(func(context *tgo.MContext) {
		hasValue = 1
	})

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)

	_, err = cn.Write([]byte{byte(tgo.MsgTypePingreq)})
	test.Nil(t, err)

	time.Sleep(time.Millisecond * 50)

	test.Equal(t, 1, hasValue)
}

func TestHandleHeartbeat(t *testing.T) {
	tg  := startTGO(t)
	tg.Use(HandleHeartbeat)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)

	_, err = cn.Write([]byte{byte(tgo.MsgTypePingreq)})
	test.Nil(t, err)
	time.Sleep(time.Millisecond * 50)


	resultBytes := make([]byte,1)
	_,err = cn.Read(resultBytes)
	test.Nil(t,err)

	test.Equal(t,tgo.MsgTypePingresp,tgo.MsgType(resultBytes[0]))
}

func TestHandleAuth(t *testing.T) {
	tg  := startTGO(t)
	tg.Use(HandleAuth)
	tg.Use(HandleHeartbeat)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)
	msgBytes := []byte{byte(tgo.MsgTypeConnect),5,0x00,0x00,0x00,0,0}
	msgBytes = append(msgBytes,[]byte("pwd")...)
	_, err = cn.Write(msgBytes)
	test.Nil(t, err)
	time.Sleep(time.Millisecond * 50)


	msg,err := tg.GetOpts().Pro.Decode(cn)
	test.Nil(t,err)

	test.Equal(t,tgo.MsgTypeConnack,msg.MsgType)
	test.Equal(t,tgo.MsgStatusAuthOk,tgo.AuthStatus(msg.VariableHeader[0]))
}

func TestHandleRevMsg(t *testing.T) {
	tg  := startTGO(t)
	tg.Use(HandleAuth)
	tg.Use(HandleHeartbeat)
	tg.Match("send",HandleRevMsg)

	cn, err := MustConnectServer(tg.Server.(*server.TCPServer).RealTCPAddr())
	test.Nil(t, err)
	msgBytes := []byte{byte(tgo.MsgTypeConnect),5,0x00,0x00,0x00,0,0}
	msgBytes = append(msgBytes,[]byte("pwd")...)
	_, err = cn.Write(msgBytes)
	test.Nil(t, err)
	msgBytes = []byte{byte(tgo.MsgTypePublish),7,0x00,0x00,0x00,0,0}
	msgBytes = append(msgBytes,[]byte("hello")...)
	_, err = cn.Write(msgBytes)
	test.Nil(t, err)
	time.Sleep(time.Millisecond * 50)

	msg,err := tg.GetOpts().Pro.Decode(cn)
	test.Nil(t,err)
	test.Equal(t,tgo.MsgTypeConnack,msg.MsgType)
	test.Equal(t,tgo.MsgStatusAuthOk,tgo.AuthStatus(msg.VariableHeader[0]))

	msg,err = tg.GetOpts().Pro.Decode(cn)
	test.Nil(t,err)
	test.Equal(t,tgo.MsgTypePuback,msg.MsgType)
	test.Equal(t,tgo.MsgStatusSuccess,tgo.MsgStatus(msg.VariableHeader[0]))
}


func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func startTGO(t *testing.T) *tgo.TGO  {
	opts := tgo.NewOptions()
	opts.TCPAddress = "0.0.0.0:0"
	opts.Log = test.NewLog(t)
	tg := tgo.New(opts)
	err :=tg.Start()
	test.Nil(t,err)
	return tg
}
