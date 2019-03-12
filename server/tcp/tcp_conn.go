package tcp

import (
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Conn struct {
	id             uint64
	conn           net.Conn
	exitChan       chan int        // Only  notify self exits
	connExitChan chan tgo.Conn // Client exit notify server
	waitGroup      tgo.WaitGroupWrapper
	sync.RWMutex
	connContextChan chan *tgo.ConnContext
	opts              atomic.Value // options
	isAuth            bool
	server tgo.Server
}

func NewConn(conn net.Conn, connContextChan chan *tgo.ConnContext, connExitChan chan tgo.Conn,server tgo.Server, opts *tgo.Options) *Conn {
	c := &Conn{
		conn:              conn,
		connContextChan: connContextChan,
		connExitChan:    connExitChan,
		exitChan:          make(chan int, 0),
		server: server,
	}
	c.storeOpts(opts)
	return c
}

func (c *Conn) storeOpts(opts *tgo.Options) {
	c.opts.Store(opts)
}

func (c *Conn) GetOpts() *tgo.Options {
	return c.opts.Load().(*tgo.Options)
}

func (c *Conn) StartIOLoop() {
	c.waitGroup.Wrap(c.ioLoop)
}

func (c *Conn) ioLoop() {
	c.Debug("开始收取消息")
	for {
		select {
		case <-c.exitChan:
			goto exit
		default:
			packet, err := c.GetOpts().Pro.DecodePacket(c.conn)
			if err != nil {
				c.Error("Decoding message failed - %v", err)
				goto exit
			}
			c.connContextChan <- tgo.NewPacketConn(packet,c,c.server)
		}
	}

exit:
	c.connExitChan <- c
	c.Debug("msgLoop is exit")
}

func (c *Conn) setDeadline(t time.Time) error {
	err := c.conn.SetDeadline(t)
	return err
}

func (c *Conn) Write(data []byte) (int,error) {
	return c.conn.Write(data)
}

func (c *Conn) Read(b []byte) (int,error) {
	return c.conn.Read(b)
}
func (c *Conn) Exit() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	close(c.exitChan)
	c.waitGroup.Wait()
	return nil
}

// --------- stateful conn -----------

func (c *Conn)  SetAuth(auth bool)  {
	c.isAuth = auth
}
func (c *Conn)  IsAuth() bool {
	return c.isAuth
}

func (c *Conn) SetId(id uint64)  {
	c.id = id
}

func (c *Conn) SetDeadline(t time.Time) error  {
	return c.conn.SetDeadline(t)
}

// --------- log -------------
func (c *Conn) Info(format string, a ...interface{}) {
	c.GetOpts().Log.Info(fmt.Sprintf("%s[%d] -> %s", c.getLogPrefix(), c.id, format), a...)
}

func (c *Conn) Error(format string, a ...interface{}) {
	c.GetOpts().Log.Error(fmt.Sprintf("%s[%d] -> %s", c.getLogPrefix(), c.id, format), a...)
}

func (c *Conn) Warn(format string, a ...interface{}) {
	c.GetOpts().Log.Warn(fmt.Sprintf("%s[%d] -> %s", c.getLogPrefix(), c.id, format), a...)
}

func (c *Conn) Debug(format string, a ...interface{}) {
	c.GetOpts().Log.Debug(fmt.Sprintf("%s[%d] -> %s", c.getLogPrefix(), c.id, format), a...)
}

func (c *Conn) Fatal(format string, a ...interface{}) {
	c.GetOpts().Log.Fatal(fmt.Sprintf("%s[%d] -> %s", c.getLogPrefix(), c.id, format), a...)
}

func (c *Conn) getLogPrefix() string {
	return "【Conn】"
}
