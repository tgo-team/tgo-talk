package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"math"
	"reflect"
	"runtime"
	"sync"
)

type HandlerFunc func(*MContext)
type AuthHandlerFunc func(MContext) error
type HandlersChain []HandlerFunc

type Route struct {
	pool            sync.Pool
	handlers        HandlersChain
	ctx             *Context
	matchHandlerMap map[string]HandlerFunc
}

func NewRoute(ctx *Context) *Route {
	r := &Route{
		handlers:        HandlersChain{},
		ctx:             ctx,
		matchHandlerMap: make(map[string]HandlerFunc, 0),
	}
	return r
}

func (r *Route) handle(context *MContext) {
	context.Next()
}

func (r *Route) Serve(context *MContext) {
	context.Ctx = r.ctx
	context.handlers = r.handlers
	r.handle(context)

	if context.Packet()!=nil && !context.IsAborted(){
		packetType := context.Packet().GetFixedHeader().PacketType
		typePath := fmt.Sprintf("type:%d",packetType)
		matchFunc := r.matchHandlerMap[typePath]
		if matchFunc!=nil {
			matchFunc(context)
		}
	}

}

func (r *Route) Use(handles ...HandlerFunc) *Route {
	r.handlers = append(r.handlers, handles...)
	return r
}

func (r *Route) Match(match string, handler HandlerFunc) {
	r.matchHandlerMap[match] = handler
}

const abortIndex int8 = math.MaxInt8 / 2

type MContext struct {
	connContext *ConnContext
	index       int8
	handlers    HandlersChain
	sync.RWMutex
	Ctx *Context
}

var pool = sync.Pool{
	New: func() interface{} {
		return allocateContext()
	},
}

func GetMContext(connContext *ConnContext) *MContext {
	mContext := pool.Get().(*MContext)
	mContext.reset()
	mContext.connContext = connContext
	return mContext
}

func allocateContext() *MContext {
	return &MContext{index: -1, handlers: nil, connContext: nil, RWMutex: sync.RWMutex{}}
}

func (c *MContext) Next() {
	c.index++
	for ; c.index < int8(len(c.handlers)); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *MContext) Packet() packets.Packet {
	return c.connContext.Packet
}

func (c *MContext) PacketType() packets.PacketType {

	return c.Packet().GetFixedHeader().PacketType
}

func (c *MContext) Conn() Conn {
	return c.connContext.Conn
}

func (c *MContext) Server() Server {

	return c.connContext.Server
}

func (c *MContext) Msg() *Msg {
	messagePacket, ok := c.connContext.Packet.(*packets.MessagePacket)
	if ok {
		msg := NewMsg(messagePacket.MessageID,messagePacket.From,messagePacket.Payload)
		msg.MessageID = messagePacket.MessageID
		msg.Payload = messagePacket.Payload
		return msg
	}
	return nil
}

func (c *MContext) Abort() {
	c.index = abortIndex
}

func (c *MContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *MContext) ReplyPacket(packet packets.Packet) error {
	data, err := c.Ctx.TGO.GetOpts().Pro.EncodePacket(packet)
	if err != nil {
		return err
	}
	_, err = c.Conn().Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *MContext) GetChannel(channelID uint64) (*Channel,error) {

	return c.Ctx.TGO.GetChannel(channelID)
}

func (c *MContext) reset() {
	c.Lock()
	defer c.Unlock()
	c.index = -1
	c.connContext = nil
	c.handlers = nil
}

func (c *MContext) current() HandlerFunc {
	if c.index < int8(len(c.handlers)) && c.index != -1 {
		return c.handlers[c.index]
	}
	return nil
}

// ---------- log --------------
func (c *MContext) Info(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("%s[%s] -> ", c.getLogPrefix(), funcName)+f, args...)
	return
}

func (c *MContext) Error(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("%s[%s] -> ", c.getLogPrefix(), funcName)+f, args...)
	return
}

func (c *MContext) Debug(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("%s[%s] -> ", c.getLogPrefix(), funcName)+f, args...)
	return
}

func (c *MContext) Warn(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("%s[%s] -> ", c.getLogPrefix(), funcName)+f, args...)
	return
}

func (c *MContext) Fatal(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("%s[%s] -> ", c.getLogPrefix(), funcName)+f, args...)
	return
}

func (c *MContext) getLogPrefix() string {
	return "【Route】"
}

func (c *MContext) currentHandleName() string {
	funcName := ""
	if c.current() != nil {
		funcName = runtime.FuncForPC(reflect.ValueOf(c.current()).Pointer()).Name()
	}

	return funcName
}
