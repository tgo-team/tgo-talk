package server

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

func init() {
	tgo.RegistryServer(func(context *tgo.Context) tgo.Server {

		return NewTCPServer(context.TGO.GetOpts())

	})
}

type TCPServer struct {
	tcpListener    net.Listener
	exitChan       chan int
	waitGroup      tgo.WaitGroupWrapper
	opts           atomic.Value // options
	cm             *clientManager
	clientExitChan chan int64 // client exit
	readMsgChan    chan *tgo.Msg
}

func NewTCPServer(opts *tgo.Options) *TCPServer {
	s := &TCPServer{
		exitChan:       make(chan int, 0),
		cm:             newClientManager(),
		clientExitChan: make(chan int64, 0),
		readMsgChan:    make(chan *tgo.Msg, 1024),
	}
	s.opts.Store(opts)
	var err error
	s.tcpListener, err = net.Listen("tcp", opts.TCPAddress)
	if err != nil {
		s.Fatal("listen (%s) failed - %s", opts.TCPAddress, err)
		os.Exit(1)
	}
	s.waitGroup.Wrap(s.clientExitLoop)
	s.waitGroup.Wrap(s.msgLoop)
	return s
}

func (s *TCPServer) storeOpts(opts *tgo.Options) {
	s.opts.Store(opts)
}

func (s *TCPServer) GetOpts() *tgo.Options {
	return s.opts.Load().(*tgo.Options)
}

func (s *TCPServer) Start() error {
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *TCPServer) ReadMsgChan() chan *tgo.Msg {

	return s.readMsgChan
}

func (s *TCPServer) WriteMsgChan() chan *tgo.Msg {

	return nil
}

func (s *TCPServer) Stop() error {
	if s.tcpListener != nil {
		err := s.tcpListener.Close()
		if err != nil {
			return err
		}
	}
	close(s.readMsgChan)
	close(s.clientExitChan)
	close(s.exitChan)
	s.waitGroup.Wait()
	s.Info("TCPServer is stopped")
	return nil
}

func (s *TCPServer) connLoop() {
	s.Info("TCP: listening on %s", s.tcpListener.Addr())
	for {
		select {
		case <-s.exitChan:
			goto exit
		default:
			cn, err := s.tcpListener.Accept()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					s.Error("temporary Accept() failure - %s", err)
					runtime.Gosched()
					continue
				}
				// theres no direct way to detect this error because it is not exposed
				if !strings.Contains(err.Error(), "use of closed network connection") {
					s.Error("listener.Accept() - %s", err)
				}
				break
			}
			s.Debug("client[%s]:connecting...", cn.RemoteAddr())
			s.waitGroup.Wrap(func() {
				s.generateClient(cn)
			})
		}
	}
exit:
	s.Info("Listener is exit")
}

func (s *TCPServer) generateClient(conn net.Conn) {

	client := newClient(conn,s.readMsgChan,s.clientExitChan,s.GetOpts())
	err := client.Start()
	if err != nil {
		s.Error("Client starts failing - %v", err)
		return
	}
	s.cm.addClient(client)

}

func (s *TCPServer) clientExitLoop() {
	for {
		select {
		case clientId := <-s.clientExitChan:
			s.Info("client[%d] is exit", clientId)
			s.cm.removeClient(clientId)
		case <-s.exitChan:
			goto exit

		}
	}
exit:
	s.Info("clientExitLoop is exit!")
}

func (s *TCPServer) msgLoop() {
	for {
		select {
		case msg := <-s.readMsgChan:
			if msg!=nil {
				s.Info("Get the message - %v", msg)
			}else{
				s.Warn("Get the message is nil")
			}
		case <-s.exitChan:
			goto exit

		}
	}
exit:
	s.Info("msgLoop is exit!")
}


func (s *TCPServer) RealTCPAddr() *net.TCPAddr {
	return s.tcpListener.Addr().(*net.TCPAddr)
}