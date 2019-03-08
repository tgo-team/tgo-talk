package tgo

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sync"
)




type HandlerFunc func( *MContext)
type AuthHandlerFunc func(MContext) error
type HandlersChain []HandlerFunc



type Route struct {
	pool     sync.Pool
	handlers HandlersChain
	ctx *Context
	matchHandlerMap map[string]HandlerFunc
}

func NewRoute(ctx *Context) *Route {
	r := &Route{
		handlers: HandlersChain{},
		ctx:ctx,
		matchHandlerMap: make(map[string]HandlerFunc,0),
	}
	return r
}

func (r *Route) handle(context *MContext) {
	context.Next()
}



func (r *Route) Serve(context *MContext) {
	context.Ctx = r.ctx
	context.Server = r.ctx.TGO.Server
	context.handlers = r.handlers
	r.handle(context)

	if context.Msg!=nil && !context.IsAborted(){
		matchFunc := r.matchHandlerMap[context.Msg.Match]
		if matchFunc!=nil {
			matchFunc(context)
		}
	}

}

func (r *Route) Use(handles ...HandlerFunc) *Route {
	r.handlers = append(r.handlers, handles...)
	return r
}


func (r *Route) Match(match string,handler HandlerFunc)  {
	r.matchHandlerMap[match] = handler
}

const abortIndex int8 = math.MaxInt8 / 2
type MContext struct {
	Msg * Msg
	index    int8
	handlers HandlersChain
	sync.RWMutex
	Ctx * Context
	Server Server
}

var pool = sync.Pool{
	New: func() interface{} {
		return allocateContext()
	},
}

func GetMContext(msg *Msg) *MContext {
	mContext := pool.Get().(*MContext)
	mContext.reset()
	mContext.Msg = msg
	return mContext
}

func allocateContext() *MContext {
	return &MContext{index: -1, handlers: nil, Msg: nil, RWMutex: sync.RWMutex{}}
}

func (c *MContext) Next() {
	c.index++
	for ; c.index < int8(len(c.handlers)); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *MContext) Abort() {
	c.index = abortIndex
}

func (c *MContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *MContext) ReplyMsg(msg *Msg) error  {
	return c.Server.SendMsg(c.Msg.ClientId,msg)
}


func (c *MContext) GetChannel(channelName string,channelType ChannelType) Channel  {

	return c.Ctx.TGO.GetChannel(channelName,channelType)
}

func (c *MContext) reset() {
	c.Lock()
	defer c.Unlock()
	c.index = -1
	c.Msg = nil
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
	c.Ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("Route[%s]:", funcName)+f, args...)
	return
}

func (c *MContext) Error(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("Route[%s]:", funcName)+f, args...)
	return
}

func (c *MContext) Debug(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("Route[%s]:", funcName)+f, args...)
	return
}

func (c *MContext) Warn(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("Route[%s]:", funcName)+f, args...)
	return
}

func (c *MContext) Fatal(f string, args ...interface{}) {
	funcName := c.currentHandleName()
	c.Ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("Route[%s]:", funcName)+f, args...)
	return
}

func (c *MContext) currentHandleName() string {
	funcName := ""
	if c.current() != nil {
		funcName = runtime.FuncForPC(reflect.ValueOf(c.current()).Pointer()).Name()
	}

	return funcName
}