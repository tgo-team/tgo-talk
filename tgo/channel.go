package tgo

import (
	"fmt"
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
	PutMsg(msg *Msg) error
	AddConsumer(id int64,consumer *Consumer)
	RemoveConsumer(id int64)
}

type PersonChannel struct {
	ctx         *Context
	channelName string
	sync.RWMutex
	consumers     map[int64]*Consumer
	messageCount  uint64
	memoryMsgChan chan *Msg
}
func NewPersonChannel(channelName string, ctx *Context) *PersonChannel {

	return &PersonChannel{
		ctx:           ctx,
		channelName:   channelName,
		consumers:     map[int64]*Consumer{},
		memoryMsgChan: make(chan *Msg, ctx.TGO.GetOpts().MemQueueSize),
	}
}

func (p *PersonChannel) PutMsg(msg *Msg) error {
	select {
	case p.memoryMsgChan <- msg:
	default:
		p.Warn("消息已满！")
	}
	atomic.AddUint64(&p.messageCount, 1)
	return nil
}

func (p *PersonChannel) AddConsumer(id int64,consumer *Consumer) {
	p.Lock()
	defer p.Unlock()
	p.consumers[id] = consumer
}

func (p *PersonChannel) RemoveConsumer(id int64) {
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
	consumers     map[int64]*Consumer
	messageCount  uint64
	memoryMsgChan chan *Msg
}

func NewGroupChannel(channelName string, ctx *Context) *GroupChannel {

	return &GroupChannel{
		ctx:           ctx,
		channelName:   channelName,
		consumers:     map[int64]*Consumer{},
		memoryMsgChan: make(chan *Msg, ctx.TGO.GetOpts().MemQueueSize),
	}
}

func (g *GroupChannel) PutMsg(msg *Msg) error {
	select {
	case g.memoryMsgChan <- msg:
	default:
		g.Warn("消息已满！")
	}
	atomic.AddUint64(&g.messageCount, 1)
	return nil
}

func (g *GroupChannel) AddConsumer(id int64,consumer *Consumer) {
	g.Lock()
	defer g.Unlock()
	g.consumers[id] = consumer
}

func (g *GroupChannel) RemoveConsumer(id int64) {
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
