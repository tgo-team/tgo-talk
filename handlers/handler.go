package handlers

import (
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"time"
)

//
//import (
//	"fmt"
//	"github.com/tgo-team/tgo-talk/tgo"
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

func HandleAuth(m *tgo.MContext)  {
	if m.Packet().GetFixedHeader().PacketType == packets.Connect {
		m.Debug("开始认证！")
		connectPacket := m.Packet().(*packets.ConnectPacket)
		if connectPacket.ClientIdentifier == 1 && string(connectPacket.Password) == "123456" {
			m.Debug("认证成功！")
			conn := m.Conn()
			if conn!=nil {
				statefulConn,ok := conn.(tgo.StatefulConn)
				if ok {
					statefulConn.SetAuth(true)
					statefulConn.SetId(connectPacket.ClientIdentifier)
					err := statefulConn.SetDeadline(time.Now().Add(m.Ctx.TGO.GetOpts().MaxHeartbeatInterval*2))
					if err!=nil {
						m.Error("SetDeadline失败 -> %v",err)
						return
					}
					statefulConn.StartIOLoop()
					err = m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodeSuccess))
					if err!=nil {
						m.Error("发送认证ACK失败 -> %v",err)
					}
					statefulServer := m.Server().(tgo.StatefulServer)
					err = statefulServer.AddConn(connectPacket.ClientIdentifier,conn)
					if err!=nil {
						m.Error("添加连接失败 -> %v",err)
						return
					}
				}
			}
		}else{
			err := m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodePasswordOrUnameError))
			if err!=nil {
				m.Error("发送认证ACK失败 -> %v",err)
			}
		}
		m.Debug("结束认证！")
	}else{
		conn := m.Conn()
		statefulConn,ok := conn.(tgo.StatefulConn)
		if ok  {
			if !statefulConn.IsAuth() {
				err := m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodeUnAuth))
				if err!=nil {
					m.Error("发送认证ACK失败 -> %v",err)
				}
				m.Abort()
				return
			}
		}else{

		}
	}
}