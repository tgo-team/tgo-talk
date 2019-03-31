package handlers

import (
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"github.com/tgo-team/tgo-talk/handlers/cmd"
	"time"
)

// HandleConnPacket 处理连接包
func HandleConnPacket(m *tgo.MContext) {

	if m.PacketType() == packets.CMD { // CMD类型不做认证判断
		return
	}

	if m.PacketType() == packets.Connect {
		m.Debug("开始认证！")
		connectPacket := m.Packet().(*packets.ConnectPacket)
		client, err := m.Storage().GetClient(connectPacket.ClientID)
		if err != nil {
			m.Error("获取客户端信息失败！-> %v", err)
			m.ReplyPacket(packets.NewConnackPacket(packets.ConnReturnCodeError))
			goto stopAuth
		}
		if client == nil {
			m.Info("客户端[%d]不存在", connectPacket.ClientID)
			m.ReplyPacket(packets.NewConnackPacket(packets.ConnReturnCodePasswordOrUnameError))
			goto stopAuth
		}
		if (connectPacket.ClientID == client.ClientID && connectPacket.Password == client.Password) || m.Ctx.TGO.GetOpts().TestOn  {
			m.Debug("认证成功！")
			conn := m.Conn()
			if conn != nil {
				statefulConn, ok := conn.(tgo.StatefulConn)
				if ok {
					statefulConn.SetAuth(true)
					statefulConn.SetID(connectPacket.ClientID)
					err := statefulConn.SetDeadline(time.Now().Add(m.Ctx.TGO.GetOpts().MaxHeartbeatInterval * 2))
					if err != nil {
						m.Error("SetDeadline失败 -> %v", err)
						goto stopAuth
					}
					statefulConn.StartIOLoop()
					m.ReplyPacket(packets.NewConnackPacket(packets.ConnReturnCodeSuccess))
					if err != nil {
						m.Error("发送认证ACK失败 -> %v", err)
						goto stopAuth
					}
					// 通知TGO认证成功
					m.Ctx.TGO.AcceptAuthenticatedChan<- tgo.NewAuthenticatedContext(connectPacket.ClientID,statefulConn)
				}
			}
		} else {
			m.Info("用户或密码不正确！")
			m.ReplyPacket(packets.NewConnackPacket(packets.ConnReturnCodePasswordOrUnameError))
			if err != nil {
				m.Error("发送认证ACK失败 -> %v", err)
				goto stopAuth
			}
		}
	stopAuth:
		m.Abort()
		m.Debug("结束认证！")
	} else {
		conn := m.Conn()
		statefulConn, ok := conn.(tgo.StatefulConn)
		if ok {
			if !statefulConn.IsAuth() {
				m.ReplyPacket(packets.NewConnackPacket(packets.ConnReturnCodeUnAuth))
				m.Abort()
				return
			}
		} else {

		}
	}
}

// HandlePingPacket 处理心跳包
func HandlePingPacket(m *tgo.MContext) {
	var err error
	statefulConn, ok := m.Conn().(tgo.StatefulConn)
	if ok {
		// 有消息往来就保活
		err = statefulConn.SetDeadline(time.Now().Add(m.Ctx.TGO.GetOpts().MaxHeartbeatInterval * 2))
		if err != nil {
			m.Error("客户端[%d]设置保活失败！-> %v", statefulConn.GetID(), err)
			return
		}
	}
	if m.PacketType() == packets.Pingreq {
		m.ReplyPacket(packets.NewPingrespPacket())
		if err != nil {
			m.Error("回复心跳包失败 -> %v", err)
			return
		}
	}
}

// HandleMessagePacket 处理消息包
func HandleMessagePacket(m *tgo.MContext) {
	messagePacket := m.Packet().(*packets.MessagePacket)
	channel, err := m.GetChannel(messagePacket.ChannelID)
	if err != nil {
		m.Error("获取Channel[%d]失败 -> %v", messagePacket.ChannelID, err)
		return
	}
	if channel != nil {
		err = channel.PutMsg(m.Msg())
		if err != nil {
			m.Error("将消息放入Channel[%d]失败！ -> %v", messagePacket.ChannelID, err)
			return
		}
		m.ReplyPacket(packets.NewMsgackPacket([]uint64{messagePacket.MessageID}))
		if err != nil {
			m.Error("回复消息[%d]的ACK失败！ -> %v", messagePacket.MessageID, err)
			return
		}
	} else {
		m.Warn("Channel[%d]不存在！", messagePacket.ChannelID)
	}
}

func HandleMsgackPacket(m *tgo.MContext) {
	msgackPacket := m.Packet().(*packets.MsgackPacket)
	m.Info("收到消息回执:%d",msgackPacket.From)
	err := m.Ctx.TGO.Storage.RemoveMsgInChannel(msgackPacket.MessageIDs,msgackPacket.From)
	if err!=nil {
		m.Error("移除消息失败！-> %v",err)
		return
	}

}

// HandleCMDPacket 处理CMD包
func HandleCMDPacket(m *tgo.MContext) {
	cmdPacket := m.Packet().(*packets.CMDPacket)
	if cmdPacket.CMD == 1 {
		cmd.Register(m)
	}
}
