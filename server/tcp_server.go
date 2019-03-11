package server

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {
	tgo.RegistryServer(func(context *tgo.Context) tgo.Server {

		return NewTCPServer(context)

	})
}

type TCPServer struct {
	tcpListener    net.Listener
	exitChan       chan int
	waitGroup      tgo.WaitGroupWrapper
	ctx           tgo.Context
	cm             *clientManager
	clientExitChan chan tgo.Client // client exit
	receivePacketChan    chan packets.Packet
	storage tgo.Storage
	opts *tgo.Options
}

func NewTCPServer(ctx *tgo.Context) *TCPServer {
	s := &TCPServer{
		exitChan:       make(chan int, 0),
		cm:             newClientManager(),
		clientExitChan: make(chan tgo.Client, 1024),
		receivePacketChan:    make(chan packets.Packet, 1024),
		opts: ctx.TGO.GetOpts(),
	}
	var err error
	s.tcpListener, err = net.Listen("tcp", s.opts.TCPAddress)
	if err != nil {
		s.Fatal("listen (%s) failed - %s", s.opts.TCPAddress, err)
		os.Exit(1)
	}
	s.waitGroup.Wrap(s.clientExitLoop)
	//s.waitGroup.Wrap(s.msgLoop)
	return s
}


func (s *TCPServer) GetOpts() *tgo.Options {
	return s.ctx.TGO.GetOpts()
}

func (s *TCPServer) Start() error {
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *TCPServer) ReceivePacketChan() chan packets.Packet {

	return s.receivePacketChan
}

func (s *TCPServer) SendMsg(to uint64, packet packets.Packet) error {
	cli := s.cm.getClient(to)
	if cli != nil {
		msgData, err := s.GetOpts().Pro.EncodePacket(packet)
		if err != nil {
			return err
		}
		return cli.Write(msgData)
	}
	return nil
}

func (s *TCPServer) SetDeadline(clientId uint64, t time.Time) error {
	cli := s.cm.getClient(clientId)
	if cli != nil {
		return cli.setDeadline(t)
	}
	return nil
}

func (s *TCPServer) Keepalive(clientId uint64) error  {
	return s.SetDeadline(clientId,time.Now().Add(s.GetOpts().MaxHeartbeatInterval*2))
}

func (s *TCPServer) GetClient(clientId uint64) tgo.Client {
	cli := s.cm.getClient(clientId)
	if cli != nil {
		return cli
	}
	return nil
}


func (s *TCPServer) AuthClient(clientId,newClientId uint64) {
	cli := s.cm.getClient(clientId)
	if cli!=nil {
		s.cm.removeClient(clientId)
		cli.id = newClientId
		cli.isAuth = true
		s.cm.addClient(newClientId,cli)
	}
}

func (s *TCPServer) ClientIsAuth(clientId uint64) bool {
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
	close(s.receivePacketChan)
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
	packet,err := s.GetOpts().Pro.DecodePacket(conn)
	if err!=nil {
		s.Error("connect is error - %v",err)
		return
	}
	if packet.GetFixedHeader().PacketType != packets.Connect {
		s.Error("连接的第一消息必须为Connect类型！")
		return
	}
	connectPacket := packet.(*packets.ConnectPacket)
	err = s.ctx.TGO.Auth.Auth(connectPacket.ClientIdentifier,string(connectPacket.Password))
	if err!=nil {
		s.Error("客户端[%d]认证失败！ - %v",connectPacket.ClientIdentifier,err)
		return
	}
	client := NewClient(conn, s.receivePacketChan, s.clientExitChan, s.GetOpts())
	if err != nil {
		s.Error("Client starts failing - %v", err)
		return
	}
	err = client.setDeadline(time.Now().Add(s.GetOpts().MaxHeartbeatInterval*2))
	if err!=nil {
		s.Error("setDeadline is error - %v",err)
		return
	}
	s.cm.addClient(connectPacket.ClientIdentifier,client)

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
