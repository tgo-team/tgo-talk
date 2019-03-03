package tgo

import "testing"

func TestRoute_Use(t *testing.T) {
	r := NewRoute(&Context{})
	r.Use(func(context *MContext) {
		if int(context.msg.MsgType) != 6 {
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
	r := NewRoute(&Context{})
	r.Use(func(context *MContext) {
		if int(context.msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})
	r.Use(func(context *MContext) {
		if int(context.msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})
	r.Match("test", func(context *MContext) {
		if int(context.msg.MsgType) != 6 {
			t.Error("消息类型错误！")
		}
	})

	r.Serve(GetMContext(&Msg{
		MsgData: MsgData{
			MsgType: 6,
		},
		Match: "test",
	}))
}
