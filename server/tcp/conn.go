package tcp

import (
	"fmt"
	"github.com/tgo-team/tgo-core/tgo"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ConnChan struct {
	connContextChan chan *tgo.ConnContext
	connExitChan    chan tgo.Conn
}

func NewConnChan(connContextChan chan *tgo.ConnContext,connExitChan    chan tgo.Conn ) *ConnChan {
	return &ConnChan{
		connContextChan:connContextChan,
		connExitChan:connExitChan,
	}
}

type Conn struct {
	id        uint64
	conn      net.Conn
	exitChan  chan int // Only  notify self exits
	waitGroup tgo.WaitGroupWrapper
	sync.RWMutex
	opts   atomic.Value // options
	isAuth bool
	ctx    *tgo.Context
	connChan *ConnChan
}

func NewConn(conn net.Conn, connChan *ConnChan, ctx *tgo.Context) *Conn {
	c := &Conn{
		conn:     conn,
		exitChan: make(chan int, 0),
		ctx:      ctx,
		connChan: connChan,
	}
	return c
}


func (c *Conn) GetOpts() *tgo.Options {
	return c.ctx.TGO.GetOpts()
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
			packet, err := c.GetOpts().Pro.DecodePacket(c)
			if err != nil {
				if err == io.EOF {
					c.Debug("正常退出。")
				}else{
					c.Error("Decoding message failed - %v", err)
				}

				goto exit
			}
			if c.connChan!=nil && c.connChan.connContextChan!=nil {
				c.connChan.connContextChan <- tgo.NewConnContext(packet, c)
			}

		}
	}

exit:
	if c.connChan!=nil && c.connChan.connExitChan!=nil {
		c.connChan.connExitChan <- c
	}
	c.Debug("msgLoop is exit")
}

func (c *Conn) setDeadline(t time.Time) error {
	err := c.conn.SetDeadline(t)
	return err
}

func (c *Conn) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *Conn) Read(b []byte) (int, error) {
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

func (c *Conn) String() string  {
	return fmt.Sprintf("id: %d",c.id)
}

// --------- stateful conn -----------

func (c *Conn) SetAuth(auth bool) {
	c.isAuth = auth
}
func (c *Conn) IsAuth() bool {
	return c.isAuth
}

func (c *Conn) SetID(id uint64) {
	c.id = id
}

func (c *Conn) GetID() uint64 {
	return c.id
}

func (c *Conn) SetDeadline(t time.Time) error {
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
