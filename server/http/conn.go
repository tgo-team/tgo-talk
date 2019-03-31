package http

import (
	"fmt"
	"github.com/tgo-team/tgo-core/tgo"
	"net/http"
	"sync"
)

type ConnChan struct {
	packetContextChan chan *tgo.PacketContext
	connExitChan    chan tgo.Conn
}

func NewConnChan(packetContextChan chan *tgo.PacketContext,connExitChan    chan tgo.Conn ) *ConnChan {
	return &ConnChan{
		packetContextChan:packetContextChan,
		connExitChan:connExitChan,
	}
}

type Conn struct {
	id        uint64
	req      *http.Request
	exitChan  chan int // Only  notify self exits
	waitGroup tgo.WaitGroupWrapper
	sync.RWMutex
	ctx    *tgo.Context
	connChan *ConnChan
	respChan chan []byte
}

func NewConn(req *http.Request,respChan chan []byte, connChan *ConnChan, ctx *tgo.Context) *Conn {
	c := &Conn{
		req:     req,
		respChan: respChan,
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
				c.Error("Decoding message failed - %v", err)
				goto exit
			}
			if c.connChan!=nil && c.connChan.packetContextChan!=nil {
				c.connChan.packetContextChan <- tgo.NewPacketContext(packet, c)
			}

		}
	}

exit:
	if c.connChan!=nil && c.connChan.connExitChan!=nil {
		c.connChan.connExitChan <- c
	}
	c.Debug("msgLoop is exit")
}


func (c *Conn) Write(data []byte) (int, error) {
	c.respChan <- data
	return 0,nil
}

func (c *Conn) Read(b []byte) (int, error) {
	return c.req.Body.Read(b)
}
func (c *Conn) Exit() error {
	if c.req != nil {
		err := c.req.Body.Close()
		if err != nil {
			return err
		}
	}
	close(c.exitChan)
	c.waitGroup.Wait()
	return nil
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

