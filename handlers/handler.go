package handlers

import (
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"time"
)

func HandleAuth(m *tgo.MContext) {
	if m.Packet().GetFixedHeader().PacketType == packets.Connect {
		m.Debug("开始认证！")
		connectPacket := m.Packet().(*packets.ConnectPacket)
		if connectPacket.ClientIdentifier == 1 && string(connectPacket.Password) == "123456" {
			m.Debug("认证成功！")
			conn := m.Conn()
			if conn != nil {
				statefulConn, ok := conn.(tgo.StatefulConn)
				if ok {
					statefulConn.SetAuth(true)
					statefulConn.SetID(connectPacket.ClientIdentifier)
					err := statefulConn.SetDeadline(time.Now().Add(m.Ctx.TGO.GetOpts().MaxHeartbeatInterval * 2))
					if err != nil {
						m.Error("SetDeadline失败 -> %v", err)
						return
					}
					statefulConn.StartIOLoop()
					err = m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodeSuccess))
					if err != nil {
						m.Error("发送认证ACK失败 -> %v", err)
					}
					m.Ctx.TGO.ConnManager.AddConn(connectPacket.ClientIdentifier, conn)
				}
			}
		} else {
			err := m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodePasswordOrUnameError))
			if err != nil {
				m.Error("发送认证ACK失败 -> %v", err)
			}
		}
		m.Debug("结束认证！")
	} else {
		conn := m.Conn()
		statefulConn, ok := conn.(tgo.StatefulConn)
		if ok {
			if !statefulConn.IsAuth() {
				err := m.ReplyMsg(packets.NewConnackPacket(packets.ConnReturnCodeUnAuth))
				if err != nil {
					m.Error("发送认证ACK失败 -> %v", err)
				}
				m.Abort()
				return
			}
		} else {

		}
	}
}

// HandleHeartbeat
func HandleHeartbeat(m *tgo.MContext) {
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
		err = m.ReplyMsg(packets.NewPingrespPacket())
		if err != nil {
			m.Error("回复心跳包失败 -> %v", err)
			return
		}
	}
}

func HandleRevMsg(m *tgo.MContext) {
	if m.PacketType() == packets.Message {
		messagePacket := m.Packet().(*packets.MessagePacket)
		channel, err := m.GetChannel(messagePacket.ChannelID)
		if err != nil {
			m.Error("获取Channel[%d]失败 -> %v", err)
			return
		}
		err = channel.PutMsg(m.Msg())
		if err != nil {
			return
		}
	}

}
