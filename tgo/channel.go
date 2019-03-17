package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
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
		err := c.Ctx.TGO.Storage.AddMsg(msgContext)
		if err != nil {
			return err
		}
	}
	atomic.AddUint64(&c.MessageCount, 1)
	return nil
}

// DeliveryMsg 投递消息
func (c *Channel) DeliveryMsg(msgCtx *MsgContext)  {
	c.Debug("开始投递消息[%d]！",msgCtx.Msg().MessageID)
	clientIDs,err := c.Ctx.TGO.Storage.GetClientIDs(msgCtx.channelID)
	if err!=nil {
		c.Error("获取管道[%d]的客户端ID集合失败！ -> %v",msgCtx.channelID,err)
		return
	}
	if clientIDs==nil || len(clientIDs)<=0 {
		c.Warn("Channel[%d]里没有客户端！",msgCtx.channelID)
		return
	}
	for _,clientID :=range clientIDs {
		if clientID == msgCtx.Msg().From { // 不发送给自己
			continue
		}
		online := IsOnline(clientID)
		if online {
			conn := c.Ctx.TGO.ConnManager.GetConn(clientID)
			if conn!=nil {
				msgPacket := packets.NewMessagePacket(msgCtx.msg.MessageID,msgCtx.channelID,msgCtx.msg.Payload)
				msgPacket.From = msgCtx.Msg().From
				msgPacketData,err := c.Ctx.TGO.GetOpts().Pro.EncodePacket(msgPacket)
				if err!=nil {
					c.Error("编码消息[%d]数据失败！-> %v",msgCtx.msg.MessageID,err)
					continue
				}
				_,err = conn.Write(msgPacketData)
				if err!=nil {
					c.Error("写入消息[%d]数据失败！-> %v",msgCtx.msg.MessageID,err)
					continue
				}
			}
		}
	}

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
