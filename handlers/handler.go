package handlers

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"time"
)

func HandleAuth(m *tgo.MContext)  {

	if m.Msg.From == 1 && string(m.Msg.Payload) == "pwd" {

	}
}

// HandleHeartbeat
func HandleHeartbeat(m *tgo.MContext)  {
	if m.Msg.MsgType == tgo.MsgTypeAuth { // Auth message is not processed
		return
	}
	var err error
	statefulServer,ok := m.Server.(tgo.StatefulServer)
	if ok {
		err = statefulServer.SetDeadline(m.Msg.ClientId,time.Now().Add(m.Ctx.TGO.GetOpts().MaxHeartbeatInterval*2))
		if err!=nil {
			m.Error("client[%d] keepalive is error - %v",m.Msg.ClientId,err)
			return
		}
	}

	if m.Msg.MsgType == tgo.MsgTypePing {
		err = m.ReplyMsg(tgo.NewPong())
		if err!=nil {
			m.Error("replyMsg is error - %v",err)
			return
		}
	}
}

func HandleRevMsg(m *tgo.MContext)  {
	
}