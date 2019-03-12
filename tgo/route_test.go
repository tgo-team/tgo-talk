package tgo

import (
	"github.com/tgo-team/tgo-talk/test"
	"testing"
)

func TestRoute_Use(t *testing.T) {
	ctx := getContext(t)
	r := NewRoute(ctx)
	r.Use(func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})

	r.Serve(GetMContext(&Msg{
		MsgData: MsgData{
			MsgType: 6,
		},
	}))
}

func TestRoute_Match(t *testing.T) {
	ctx := getContext(t)
	r := NewRoute(ctx)
	var pass bool
	r.Use(func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})
	r.Use(func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})
	r.Match("test", func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
		pass = true
	})

	r.Serve(GetMContext(&Msg{
		MsgData: MsgData{
			MsgType: 6,
		},
		Match: "test",
	}))

	test.Equal(t,true,pass)
}

func TestRoute_Abort(t *testing.T) {
	ctx := getContext(t)
	r := NewRoute(ctx)
	var abort bool = true
	r.Use(func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})
	r.Use(func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
		context.Abort()
	})
	r.Match("test", func(context *MContext) {
		if int(context.Msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
		abort = false
	})

	r.Serve(GetMContext(&Msg{
		MsgData: MsgData{
			MsgType: 6,
		},
		Match: "test",
	}))
	test.Equal(t,true,abort)
}

func getContext(t *testing.T) *Context  {
	opts := NewOptions()
	opts.Log = test.NewLog(t)
	RegistryServer(func(context *Context) Server {
		return &ServerTest{}
	})
	RegistryStorage(func(context *Context) Storage {
		return  &StorageTest{}
	})
	tg := New(opts)
	return &Context{TGO:tg}
}

type ServerTest struct {

}


func (s *ServerTest) Start() error {
	return nil
}
func (s *ServerTest)  ReceiveMsgChan() chan *Msg {
	return nil
}
func (s *ServerTest)  SendMsg(to int64,msg *Msg) error {
	return nil
}
func (s *ServerTest)  Stop() error {
	return nil
}

type StorageTest struct {

}

func (s *StorageTest) SaveMsg(msg *Msg) error {
	return nil
}
func (s *StorageTest) ReceiveMsgChan() chan *Msg {
	return nil
}