package server

import (
	"github.com/tgo-team/tgo-chat/protocol"
	"github.com/tgo-team/tgo-chat/test"
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"testing"
)

func TestClient_StartAndStop(t *testing.T) {
	serverConn,client,readMsgChan,exitChan := getClient(t)
	err := client.Start()
	test.Nil(t, err)
	go func() {
		msg := <-readMsgChan
		test.Equal(t, 6, int(msg.MsgType))
		err = client.Stop()
		test.Nil(t,err)
	}()

	_, err = serverConn.Write([]byte{0x06})
	test.Nil(t, err)

	<-exitChan
}

func TestClientManager_addClient(t *testing.T)  {
	_,client,_,_ := getClient(t)
	err := client.Start()
	test.Nil(t,err)
	cm := newClientManager()
	cm.addClient(client)
	test.Equal(t, 1,len(cm.clients))
}

func TestClientManager_removeClient(t *testing.T)  {
	_,client,_,_ := getClient(t)
	err := client.Start()
	test.Nil(t,err)
	cm := newClientManager()
	clientId := cm.addClient(client)

	cm.removeClient(clientId)

	test.Equal(t,0, len(cm.clients))
}

func getClient(t *testing.T) (net.Conn,*client,chan *tgo.Msg,chan int64)  {
	serverConn, clientConn := net.Pipe()
	readMsgChan := make(chan *tgo.Msg, 100)
	exitChan := make(chan int64, 0)
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)
	opts.Pro = protocol.NewTGO()
	client := newClient(clientConn, readMsgChan, exitChan, opts)

	return serverConn,client,readMsgChan,exitChan
}