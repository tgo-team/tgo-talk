package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/pqueue"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ------------ channel ----------------
type ChannelType int

const (
	ChannelTypePerson ChannelType = iota // 个人管道
	ChannelTypeGroup                     // 群组管道
)

type Channel struct {
	ChannelID    uint64
	ChannelType  ChannelType
	MessageCount uint64
	sync.RWMutex
	Ctx *Context

	inFlightMessages map[uint64]*Msg
	inFlightPQ       pqueue.PriorityQueue
	inFlightMutex    sync.Mutex

	connMap map[uint64]*Conn
}

func NewChannel(channelID uint64, channelType ChannelType, ctx *Context) *Channel {
	c := &Channel{
		connMap:     map[uint64]*Conn{},
		ChannelID:   channelID,
		ChannelType: channelType,
		Ctx:         ctx,
	}
	//c.initPQ()
	return c
}

func (c *Channel) initPQ() {
	pqSize := int(math.Max(1, float64(c.Ctx.TGO.GetOpts().MemQueueSize)/10))

	c.inFlightMutex.Lock()
	c.inFlightMessages = make(map[uint64]*Msg)
	c.inFlightPQ = pqueue.New(pqSize)
	c.inFlightMutex.Unlock()
}

func (c *Channel) PutMsg(msg *Msg) error {
	msgContext := NewMsgContext(msg, c.ChannelID)
	select {
	case c.Ctx.TGO.memoryMsgChan <- msgContext:
	default:
		c.Warn("内存消息已满，进入持久化存储！")
		err := c.Ctx.TGO.Storage.SaveMsg(msgContext)
		if err != nil {
			return err
		}
	}
	atomic.AddUint64(&c.MessageCount, 1)
	return nil
}

func (c *Channel) StartInFlightTimeout(msg *Msg, clientID int64, timeout time.Duration) error {
	now := time.Now()
	item := &pqueue.Item{Value: msg, Priority: now.Add(timeout).UnixNano()}
	c.addToInFlightPQ(item)
	return nil
}

func (c *Channel) addToInFlightPQ(item *pqueue.Item) {
	c.inFlightMutex.Lock()
	c.inFlightPQ.Push(item)
	c.inFlightMutex.Unlock()
}

func (c *Channel) String() string {
	return fmt.Sprintf("ChannelID: %d ChannelType: %d MessageCount: %d", c.ChannelID, c.ChannelType, c.MessageCount)
}

// ---------- log --------------

func (c *Channel) Info(f string, args ...interface{}) {
	c.Ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("%s[%d] -> ", c.getLogPrefix(), c.ChannelID)+f, args...)
	return
}

func (c *Channel) Error(f string, args ...interface{}) {
	c.Ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("%s[%d] -> ", c.getLogPrefix(), c.ChannelID)+f, args...)
	return
}

func (c *Channel) Debug(f string, args ...interface{}) {
	c.Ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("%s[%d] -> ", c.getLogPrefix(), c.ChannelID)+f, args...)
	return
}

func (c *Channel) Warn(f string, args ...interface{}) {
	c.Ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("%s[%d] -> ", c.getLogPrefix(), c.ChannelID)+f, args...)
	return
}

func (c *Channel) Fatal(f string, args ...interface{}) {
	c.Ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("%s[%d] -> ", c.getLogPrefix(), c.ChannelID)+f, args...)
	return
}

func (c *Channel) getLogPrefix() string {
	return "【Chanel】"
}
