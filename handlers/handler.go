package handlers

import (
	"fmt"
	"github.com/tgo-team/tgo-chat/tgo"
)

func HandleAuth(m *tgo.MContext)  {
	if string(m.Msg.Payload) == "pwd" {
		statefulServer,ok := m.Server.(tgo.StatefulServer)
		if ok {
			// 如果是StatefulServer client需要设置为认证 并且更新clientId
			statefulServer.AuthClient(m.Msg.ClientId,m.Msg.UID)
			m.Msg.ClientId = m.Msg.UID
		}
		// 获取认证用户的设备ID 并加入到自己的channel里
		channel := m.GetChannel(fmt.Sprintf("%d",m.Msg.UID))
		channel.AddConsumer(m.Msg.UID)

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
		// 如果是StatefulServer 由消息往来就保活
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

	//m.Server.SendMsg()
}