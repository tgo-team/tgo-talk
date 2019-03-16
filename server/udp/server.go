package udp

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo"
	"net"
)

func init() {
	tgo.RegistryServer(func(context *tgo.Context) tgo.Server {
		return NewServer(context)

	})
}

type Server struct {
	exitChan        chan int
	waitGroup       tgo.WaitGroupWrapper
	connExitChan    chan tgo.Conn // client exit
	connContextChan chan *tgo.ConnContext
	storage         tgo.Storage
	opts            *tgo.Options
	pro             tgo.Protocol
	ctx             *tgo.Context
	conn            *net.UDPConn
	udpAddr         *net.UDPAddr
}

func NewServer(ctx *tgo.Context) *Server {
	s := &Server{
		exitChan:        make(chan int, 0),
		connExitChan:    make(chan tgo.Conn, 1024),
		connContextChan: ctx.TGO.ConnContextChan,
		opts:            ctx.TGO.GetOpts(),
		pro:             ctx.TGO.GetOpts().Pro,
		ctx:             ctx,
	}
	var err error
	s.udpAddr, err = net.ResolveUDPAddr("udp", ctx.TGO.GetOpts().UDPAddress)
	if err != nil {
		panic(err)
	}
	s.conn, err = net.ListenUDP("udp", s.udpAddr)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Server) Start() error {

	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) connLoop() {
	s.Info("开始监听 -> %s", s.udpAddr.String())
	data := make([]byte,0x7fff)
	for  {
		n,addr,err := s.conn.ReadFromUDP(data)
		if err!=nil {
			panic(err)
		}

		packet, err := s.pro.DecodePacket(bytes.NewBuffer(data[:n]))
		if err != nil {
			s.Error("解析连接数据失败！-> %v", err)
			s.exitChan <- 1
			return
		}
		cn := NewConn(s.conn,addr,NewConnChan(s.connContextChan,s.connExitChan),s.ctx)
		s.connContextChan <- tgo.NewConnContext(packet,cn)
	}


}

// --------- log -------------
func (s *Server) Info(format string, a ...interface{}) {
	s.opts.Log.Info(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Error(format string, a ...interface{}) {
	s.opts.Log.Error(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Warn(format string, a ...interface{}) {
	s.opts.Log.Warn(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Debug(format string, a ...interface{}) {
	s.opts.Log.Debug(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Fatal(format string, a ...interface{}) {
	s.opts.Log.Fatal(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) getLogPrefix() string {
	return "UDPServer"
}
