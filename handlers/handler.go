package handlers
//
//import (
//	"fmt"
//	"github.com/tgo-team/tgo-chat/tgo"
//)
//
//func HandleAuth(m *tgo.MContext)  {
//	errorAuthACK := tgo.NewAuthACK(0)
//	statefulServer,ok := m.Server.(tgo.StatefulServer)
//	if ok {
//		if statefulServer.ClientIsAuth(m.Msg.ClientId) {
//			return
//		}else{
//			if m.Msg.MsgType != tgo.MsgTypeConnect {
//				err := m.ReplyMsg(errorAuthACK)
//				if err!=nil {
//					m.Error("ReplyMsg is error - %v",err)
//					return
//				}
//				m.Abort()
//				return
//			}
//		}
//	}
//	if string(m.Msg.Payload) == "pwd" {
//		if ok {
//			// 如果是StatefulServer client需要设置为认证 并且更新clientId
//			statefulServer.AuthClient(m.Msg.ClientId,m.Msg.UID)
//			m.Msg.ClientId = m.Msg.UID
//		}
//		// 获取认证用户的设备ID 并加入到自己的channel里
//		channel := m.GetChannel(fmt.Sprintf("%d",m.Msg.UID),tgo.ChannelTypePerson)
//		var deviceId int64= 0
//		channel.AddConsumer(fmt.Sprintf("%d-%d",m.Msg.UID,deviceId),tgo.NewConsumer(m.Msg.UID,deviceId))
//
//		err := m.ReplyMsg(tgo.NewAuthACK(tgo.MsgStatusAuthOk))
//		if err!=nil {
//			m.Error("ReplyMsg is error - %v",err)
//			return
//		}
//	}else{
//		err := m.ReplyMsg(errorAuthACK)
//		if err!=nil {
//			m.Error("ReplyMsg is error - %v",err)
//			return
//		}
//		m.Abort()
//	}
//}
//// HandleHeartbeat
//func HandleHeartbeat(m *tgo.MContext)  {
//	var err error
//	statefulServer,ok := m.Server.(tgo.StatefulServer)
//	if ok {
//		// 如果是StatefulServer 由消息往来就保活
//		err = statefulServer.Keepalive(m.Msg.ClientId)
//		if err!=nil {
//			m.Error("client[%d] keepalive is error - %v",m.Msg.ClientId,err)
//			return
//		}
//	}
//	if m.Msg.MsgType == tgo.MsgTypePingreq {
//		err = m.ReplyMsg(tgo.NewPong())
//		if err!=nil {
//			m.Error("replyMsg is error - %v",err)
//			return
//		}
//	}
//}
//
//func HandleRevMsg(m *tgo.MContext)  {
//
//
//	channel := m.GetChannel(fmt.Sprintf("%d",m.Msg.UID),tgo.ChannelTypePerson)
//	err := channel.PutMsg(m.Msg)
//	if err!=nil {
//		m.Error("PutMsg is error - %v",err)
//		return
//	}
//
//	toChannel := m.GetChannel(fmt.Sprintf("%d",m.Msg.ToUID),tgo.ChannelTypePerson)
//	err  = toChannel.PutMsg(m.Msg)
//	if err!=nil {
//		m.Error("PutMsg is error - %v",err)
//		return
//	}
//	m.ReplyMsg(tgo.NewSendMsgACK(2334,tgo.MsgStatusSuccess))
//	//m.Server.SendMsg()
//}