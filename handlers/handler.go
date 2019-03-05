package handlers

import (
	"github.com/tgo-team/tgo-chat/tgo"
)

func HandleAuth(m *tgo.MContext)  {
	if string(m.Msg.Payload) == "pwd" {
		statefulServer,ok := m.Server.(tgo.StatefulServer)
		if ok {
			statefulServer.AuthClient(m.Msg.ClientId,m.Msg.From)
			m.Msg.ClientId = m.Msg.From
		}
		err := m.ReplyMsg(tgo.NewAuthAck(tgo.MsgStatusAuthOk))
		if err!=nil {
			m.Error("ReplyMsg is error - %v",err)
			return
		}
	}else{
		err := m.ReplyMsg(tgo.NewAuthAck(0))
		if err!=nil {
			m.Error("ReplyMsg is error - %v",err)
			return
		}
		m.Abort()
	}
}
// HandleHeartbeat
func HandleHeartbeat(m *tgo.MContext)  {
	var err error
	statefulServer,ok := m.Server.(tgo.StatefulServer)
	if ok {
		err = statefulServer.Keepalive(m.Msg.ClientId)
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

	m.Server.SendMsg()
}