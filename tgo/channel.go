package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"sync"
	"sync/atomic"
)

type Consumer struct {
	uid      int64
	deviceId int64
}

func NewConsumer(uid int64, deviceId int64) *Consumer {
	return &Consumer{uid: uid, deviceId: deviceId}
}

type Channel interface {
	PutPacket(packet packets.Packet) error
	AddConsumer(id string,consumer *Consumer)
	RemoveConsumer(id string)
}

type ChannelType int
const (
	ChannelTypePerson ChannelType = iota
	ChannelTypeGroup
)


// ------------ person channel ----------------

type PersonChannel struct {
	ctx         *Context
	channelName string
	sync.RWMutex
	consumers     map[string]*Consumer
	messageCount  uint64
}
func NewPersonChannel(channelName string, ctx *Context) *PersonChannel {

	return &PersonChannel{
		ctx:           ctx,
		channelName:   channelName,
		consumers:     map[string]*Consumer{},
	}
}

func (p *PersonChannel) PutPacket(packet packets.Packet) error {
	select {
	case p.ctx.TGO.memoryPacketChan <- packet:
	default:
		p.Warn("内存消息已满，进入持久化存储！")
		err := p.ctx.TGO.Storage.SaveMsg(packet)
		if err!=nil {
			return err
		}
	}
	atomic.AddUint64(&p.messageCount, 1)
	return nil
}

func (p *PersonChannel) AddConsumer(id string,consumer *Consumer) {
	p.Lock()
	defer p.Unlock()
	p.consumers[id] = consumer
}

func (p *PersonChannel) RemoveConsumer(id string) {
	p.Lock()
	defer p.Unlock()

	_, ok := p.consumers[id]
	if !ok {
		return
	}
	delete(p.consumers, id)
}

// ------------ group channel ----------------

type GroupChannel struct {
	ctx         *Context
	channelName string
	sync.RWMutex
	consumers     map[string]*Consumer
	messageCount  uint64
}

func NewGroupChannel(channelName string, ctx *Context) *GroupChannel {

	return &GroupChannel{
		ctx:           ctx,
		channelName:   channelName,
		consumers:     map[string]*Consumer{},
	}
}

func (g *GroupChannel) PutPacket(packet packets.Packet) error {
	select {
	case g.ctx.TGO.memoryPacketChan <- packet:
	default:
		g.Warn("内存消息已满，进入持久化存储！")
		err := g.ctx.TGO.Storage.SaveMsg(packet)
		if err!=nil {
			return err
		}
	}
	atomic.AddUint64(&g.messageCount, 1)
	return nil
}

func (g *GroupChannel) AddConsumer(id string,consumer *Consumer) {
	g.Lock()
	defer g.Unlock()
	g.consumers[id] = consumer
}

func (g *GroupChannel) RemoveConsumer(id string) {
	g.Lock()
	defer g.Unlock()

	_, ok := g.consumers[id]
	if !ok {
		return
	}
	delete(g.consumers, id)
}

// ---------- log --------------
func (g *GroupChannel) Info(f string, args ...interface{}) {
	g.ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("GroupChannel[%s]:", g.channelName)+f, args...)
	return
}

func (g *GroupChannel) Error(f string, args ...interface{}) {
	g.ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("GroupChannel[%s]:", g.channelName)+f, args...)
	return
}

func (g *GroupChannel) Debug(f string, args ...interface{}) {
	g.ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("GroupChannel[%s]:", g.channelName)+f, args...)
	return
}

func (g *GroupChannel) Warn(f string, args ...interface{}) {
	g.ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("GroupChannel[%s]:", g.channelName)+f, args...)
	return
}

func (g *GroupChannel) Fatal(f string, args ...interface{}) {
	g.ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("GroupChannel[%s]:", g.channelName)+f, args...)
	return
}

func (p *PersonChannel) Info(f string, args ...interface{}) {
	p.ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("GroupChannel[%s]:", p.channelName)+f, args...)
	return
}

func (p *PersonChannel) Error(f string, args ...interface{}) {
	p.ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("GroupChannel[%s]:", p.channelName)+f, args...)
	return
}

func (p *PersonChannel) Debug(f string, args ...interface{}) {
	p.ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("GroupChannel[%s]:", p.channelName)+f, args...)
	return
}

func (p *PersonChannel) Warn(f string, args ...interface{}) {
	p.ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("GroupChannel[%s]:", p.channelName)+f, args...)
	return
}

func (p *PersonChannel) Fatal(f string, args ...interface{}) {
	p.ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("GroupChannel[%s]:", p.channelName)+f, args...)
	return
}
