package cmd

import (
	"bytes"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
)

func Register(m *tgo.MContext) {
	m.Info("Register")
	cmdPacket := m.Packet().(*packets.CMDPacket)

	payloadReader := bytes.NewBuffer(cmdPacket.Payload)

	clientID := packets.DecodeUint64(payloadReader)
	password := packets.DecodeString(payloadReader)

	client, err := m.Storage().GetClient(clientID)
	if err != nil {
		m.Error("查询客户端[%d]失败！ -> %v", clientID, err)
		replyCMDPacketError(m, CMDRegisterAck, RegisterError)
		return
	}
	if client != nil {
		m.Error("客户端[%d]已存在！", clientID)
		replyCMDPacketError(m, CMDRegisterAck, RegisterClientExist)
		return
	}
	err = m.Storage().AddClient(tgo.NewClient(clientID, password))
	if err != nil {
		replyCMDPacketError(m, CMDRegisterAck, RegisterError)
		return
	}
	var channelID = clientID
	// 添加个人管道
	err = m.Storage().AddChannel(tgo.NewChannel(channelID,tgo.ChannelTypePerson,m.Ctx))
	if err!=nil {
		m.Error("添加Channel失败！-> %v",err)
		replyCMDPacketError(m, CMDRegisterAck, RegisterError)
		return
	}

	if err := m.Storage().Bind(clientID,channelID);err!=nil {
		m.Error("绑定Channel失败！-> %v",err)
		replyCMDPacketError(m, CMDRegisterAck, RegisterError)
		return
	}
	replyCMDPacketSuccess(m, CMDRegisterAck)

}

func replyCMDPacketError(m *tgo.MContext, cmd uint16, status uint16) {
	m.ReplyPacket(packets.NewCMDPacket(cmd, packets.EncodeUint16(status)))
}

func replyCMDPacketSuccess(m *tgo.MContext, cmd uint16) {
	m.ReplyPacket(packets.NewCMDPacket(cmd, packets.EncodeUint16(uint16(CMDSuccess))))
}
