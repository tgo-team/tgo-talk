package server

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
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
	clientExitChan chan tgo.Client // client exit
	readMsgChan    chan *tgo.Msg
}

func NewTCPServer(opts *tgo.Options) *TCPServer {
	s := &TCPServer{
		exitChan:       make(chan int, 0),
		cm:             newClientManager(),
		clientExitChan: make(chan tgo.Client, 1024),
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
	//s.waitGroup.Wrap(s.msgLoop)
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

func (s *TCPServer) SendMsg(to int64, msg *tgo.Msg) error {
	cli := s.cm.getClient(to)
	if cli != nil {
		msgData, err := s.GetOpts().Pro.Encode(msg)
		if err != nil {
			return err
		}
		return cli.Write(msgData)
	}
	return nil
}

func (s *TCPServer) SetDeadline(clientId int64, t time.Time) error {
	cli := s.cm.getClient(clientId)
	if cli != nil {
		return cli.setDeadline(t)
	}
	return nil
}

func (s *TCPServer) Keepalive(clientId int64) error  {
	return s.SetDeadline(clientId,time.Now().Add(s.GetOpts().MaxHeartbeatInterval*2))
}

func (s *TCPServer) GetClient(clientId int64) tgo.Client {
	cli := s.cm.getClient(clientId)
	if cli != nil {
		return cli
	}
	return nil
}


func (s *TCPServer) AuthClient(clientId,newClientId int64) {
	cli := s.cm.getClient(clientId)
	if cli!=nil {
		s.cm.removeClient(clientId)
		cli.id = newClientId
		cli.isAuth = true
		s.cm.addClient(newClientId,cli)
	}
}

func (s *TCPServer) ClientIsAuth(clientId int64) bool {
	cli := s.cm.getClient(clientId)
	if cli!=nil {
		return cli.isAuth
	}
	return false
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
	err := conn.SetDeadline(time.Now().Add(time.Second * 2)) // 2 seconds authentication time
	if err != nil {
		s.Error("SetDeadline is error - %v", err)
		return
	}
	client := NewClient(conn, s.readMsgChan, s.clientExitChan, s.GetOpts())
	if err != nil {
		s.Error("Client starts failing - %v", err)
		return
	}
	clientId := atomic.AddInt64(&s.cm.clientIDSequence, -1)
	s.cm.addClient(clientId,client)

}

func (s *TCPServer) clientExitLoop() {
	for {
		select {
		case cli := <-s.clientExitChan:
			s.Info("client[%v] is exit", cli)
			//s.cm.removeClient(cli)
		case <-s.exitChan:
			goto exit

		}
	}
exit:
	s.Info("clientExitLoop is exit!")
}

func (s *TCPServer) RealTCPAddr() *net.TCPAddr {
	return s.tcpListener.Addr().(*net.TCPAddr)
}
