package server

import (
	"fmt"
	"github.com/tgo-team/tgo-chat/tgo"
	"net"
	"sync"
	"sync/atomic"
)

type client struct {
	id             int64
	authId         int64
	conn           net.Conn
	exitChan       chan int   // Only  notify self exits
	clientExitChan chan int64 // Client exit notify server
	waitGroup      tgo.WaitGroupWrapper
	sync.RWMutex
	readMsgChan chan *tgo.Msg
	opts        atomic.Value // options
}

func newClient(conn net.Conn, readMsgChan chan *tgo.Msg, clientExitChan chan int64, opts *tgo.Options) *client {
	c := &client{
		conn:           conn,
		readMsgChan:    readMsgChan,
		clientExitChan: clientExitChan,
		exitChan:       make(chan int, 0),
	}
	c.storeOpts(opts)
	return c
}

func (c *client) Start() error {
	c.waitGroup.Wrap(c.msgLoop)
	return nil
}

func (c *client) Stop() error {
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.exitChan)
	c.waitGroup.Wait()
	return nil
}

func (c *client) storeOpts(opts *tgo.Options) {
	c.opts.Store(opts)
}

func (c *client) GetOpts() *tgo.Options {
	return c.opts.Load().(*tgo.Options)
}

func (c *client) msgLoop() {
	for {
		select {
		case <-c.exitChan:
			goto exit
		default:
			msg, err := c.GetOpts().Pro.Decode(c.conn)
			if err != nil {
				c.Warn("Decoding message failed - %v", err)
				goto exit
			}
			c.readMsgChan <- msg
		}
	}

exit:
	c.clientExitChan<-1
	c.Info("msgLoop is exit")
}

type clientManager struct {
	clients          map[int64]*client
	clientLock       sync.RWMutex
	clientIDSequence int64
}

func newClientManager() *clientManager {

	return &clientManager{
		clients: make(map[int64]*client),
	}
}

func (cm *clientManager) addClient(client *client) int64 {
	clientId := atomic.AddInt64(&cm.clientIDSequence, 1)
	client.id = clientId
	cm.clientLock.Lock()
	cm.clients[client.id] = client
	cm.clientLock.Unlock()
	return client.id
}

func (cm *clientManager) removeClient(clientId int64) {
	cm.clientLock.Lock()
	_, ok := cm.clients[clientId]
	if !ok {
		cm.clientLock.Unlock()
		return
	}
	delete(cm.clients, clientId)
	cm.clientLock.Unlock()
}

// --------- log -------------
func (c *client) Info(format string, a ...interface{}) {
	c.GetOpts().Log.Info(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *client) Error(format string, a ...interface{}) {
	c.GetOpts().Log.Error(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *client) Warn(format string, a ...interface{}) {
	c.GetOpts().Log.Warn(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *client) Debug(format string, a ...interface{}) {
	c.GetOpts().Log.Debug(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *client) Fatal(format string, a ...interface{}) {
	c.GetOpts().Log.Fatal(fmt.Sprintf("【%s】%s", c.getLogPrefix(), format), a...)
}

func (c *client) getLogPrefix() string {
	return "Client"
}
