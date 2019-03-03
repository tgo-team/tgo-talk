package server

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"os"
	"runtime"
	"strings"
)

func init() {
	tgo.RegistryServer(func(context *tgo.Context) tgo.Server {

		return NewTCPServer(context.TGO.GetOpts())

	})
}

type TCPServer struct {
	opts        *tgo.Options
	tcpListener net.Listener
	exitChan    chan int
	waitGroup   tgo.WaitGroupWrapper
}

func NewTCPServer(opts *tgo.Options) *TCPServer {
	s := &TCPServer{
		opts:     opts,
		exitChan: make(chan int, 0),
	}
	var err error
	s.tcpListener, err = net.Listen("tcp", opts.TCPAddress)
	if err != nil {
		s.Fatal("listen (%s) failed - %s", opts.TCPAddress, err)
		os.Exit(1)
	}
	return s
}

func (s *TCPServer) Start() error {
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *TCPServer) MsgChan() chan tgo.Msg {

	return nil
}

func (s *TCPServer) Stop() error {
	if s.tcpListener != nil {
		err := s.tcpListener.Close()
		if err != nil {
			return err
		}
	}
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
		}
	}
exit:
	s.Info("Listener is exit")
}

func msgLoop(conn net.Conn)  {

}
