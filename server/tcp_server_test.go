package server

import (
	"github.com/tgo-team/tgo-chat/test"
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"testing"
	"time"
)

func TestTCPServer_Start(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	s,_ := MustStartServer(opts)
	err := s.Start()
	test.Nil(t,err)
	time.Sleep(50 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestTCPServer_ReadMsgChan(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	s,tcpAddr := MustStartServer(opts)
	err := s.Start()
	test.Nil(t,err)

	readMsgChan := s.ReceiveMsgChan()

	go func() {
		msg :=<-readMsgChan
		test.Equal(t,tgo.MsgTypePing,msg.MsgType)
	}()

	conn,err := MustConnectServer(tcpAddr)
	test.Nil(t,err)

	_,err = conn.Write([]byte{0x06})
	test.Nil(t,err)

	time.Sleep(time.Millisecond*50)
}


func MustStartServer(opts *tgo.Options) (*TCPServer, *net.TCPAddr) {
	opts.TCPAddress = "127.0.0.1:0"
	opts.HTTPAddress = "127.0.0.1:0"
	opts.HTTPSAddress = "127.0.0.1:0"
	s := NewTCPServer(opts)
	return s, s.RealTCPAddr()
}


func MustConnectServer(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}