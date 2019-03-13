package tgo

import (
	"fmt"
	"sync"
	"sync/atomic"
)


// ------------ person channel ----------------

type Channel struct {
	ctx         *Context
	channelID uint64
	sync.RWMutex
	messageCount  uint64
}
func NewChannel(channelID uint64, ctx *Context) *Channel {

	return &Channel{
		ctx:           ctx,
		channelID:   channelID,
	}
}

func (c *Channel) PutMsg(msg *Msg) error {
	msgContext := NewMsgContext(msg,c.channelID)
	select {
	case c.ctx.TGO.memoryMsgChan <- msgContext:
	default:
		c.Warn("内存消息已满，进入持久化存储！")
		err := c.ctx.TGO.Storage.SaveMsg(msgContext)
		if err!=nil {
			return err
		}
	}
	atomic.AddUint64(&c.messageCount, 1)
	return nil
}



// ---------- log --------------

func (c *Channel) Info(f string, args ...interface{}) {
	c.ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("%s[%d] -> ",c.getLogPrefix(), c.channelID)+f, args...)
	return
}

func (c *Channel) Error(f string, args ...interface{}) {
	c.ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("%s[%d] -> ",c.getLogPrefix(), c.channelID)+f, args...)
	return
}

func (c *Channel) Debug(f string, args ...interface{}) {
	c.ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("%s[%d] -> ",c.getLogPrefix(), c.channelID)+f, args...)
	return
}

func (c *Channel) Warn(f string, args ...interface{}) {
	c.ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("%s[%d] -> ",c.getLogPrefix(), c.channelID)+f, args...)
	return
}

func (c *Channel) Fatal(f string, args ...interface{}) {
	c.ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("%s[%d] -> ",c.getLogPrefix(), c.channelID)+f, args...)
	return
}

func (c *Channel) getLogPrefix() string {
	return "【Chanel】"
}
