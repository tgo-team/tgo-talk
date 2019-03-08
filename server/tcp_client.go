package server

import (
	"fmt"
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	id             int64
	conn           net.Conn
	exitChan       chan int   // Only  notify self exits
	clientExitChan chan tgo.Client // Client exit notify server
	waitGroup      tgo.WaitGroupWrapper
	sync.RWMutex
	receiveMsgChan chan *tgo.Msg
	opts        atomic.Value // options
	isAuth bool
}

func NewClient(conn net.Conn, receiveMsgChan chan *tgo.Msg, clientExitChan chan tgo.Client, opts *tgo.Options) *Client {
	c := &Client{
		conn:           conn,
		receiveMsgChan:    receiveMsgChan,
		clientExitChan: clientExitChan,
		exitChan:       make(chan int, 0),
	}
	c.storeOpts(opts)
	c.waitGroup.Wrap(c.msgLoop)
	return c
}


func (c *Client) storeOpts(opts *tgo.Options) {
	c.opts.Store(opts)
}

func (c *Client) GetOpts() *tgo.Options {
	return c.opts.Load().(*tgo.Options)
}

func (c *Client) msgLoop() {
	for {
		select {
		case <-c.exitChan:
			goto exit
		default:
			msg, err := c.GetOpts().Pro.Decode(c.conn)
			if err != nil {
				c.Error("Decoding message failed - %v", err)
				goto exit
			}
			msg.ClientId = c.id
			msg.UID = c.id
			c.GetOpts().Monitor.TraceMsg("received",msg.Id)
			c.receiveMsgChan <- msg
		}
	}

exit:
	c.clientExitChan<-c
	c.Info("msgLoop is exit")
}


func (c *Client) setDeadline(t time.Time) error  {
	err := c.conn.SetDeadline(t)
	return err
}

func (c *Client) Write(data []byte) error {
	_,err := c.conn.Write(data)
	return err
}
func (c *Client) Exit() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err!=nil {
			return err
		}
	}
	close(c.exitChan)
	c.waitGroup.Wait()
	return nil
}





// --------- log -------------
func (c *Client) Info(format string, a ...interface{}) {
	c.GetOpts().Log.Info(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *Client) Error(format string, a ...interface{}) {
	c.GetOpts().Log.Error(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *Client) Warn(format string, a ...interface{}) {
	c.GetOpts().Log.Warn(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *Client) Debug(format string, a ...interface{}) {
	c.GetOpts().Log.Debug(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *Client) Fatal(format string, a ...interface{}) {
	c.GetOpts().Log.Fatal(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *Client) getLogPrefix() string {
	return "Client"
}
