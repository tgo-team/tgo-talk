package tgo

import (
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
	context.handlers = r.handlers
	r.handle(context)

	if context.msg!=nil{
		matchFunc := r.matchHandlerMap[context.msg.Match]
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

type MContext struct {
	msg * Msg
	index    int
	handlers HandlersChain
	sync.RWMutex
	ctx * Context
}

var pool = sync.Pool{
	New: func() interface{} {
		return allocateContext()
	},
}

func GetMContext(msg *Msg) *MContext {
	mContext := pool.Get().(*MContext)
	mContext.Reset()
	mContext.msg = msg
	return mContext
}

func allocateContext() *MContext {
	return &MContext{index: -1, handlers: nil, msg: nil, RWMutex: sync.RWMutex{}}
}

func (c *MContext) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *MContext) Current() HandlerFunc {
	if c.index < len(c.handlers) && c.index != -1 {
		return c.handlers[c.index]
	}
	return nil
}

func (c *MContext) Abort() {
	c.index = len(c.handlers)
}

func (c *MContext) IsAborted() bool {
	return c.index >= len(c.handlers)
}

func (c *MContext) Reset() {
	c.Lock()
	defer c.Unlock()
	c.index = -1
	c.msg = nil
	c.handlers = nil

}