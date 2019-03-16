package cmd

import (
	"bytes"
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

func Register(m *tgo.MContext) {
	m.Info("Register")
	cmdPacket := m.Packet().(*packets.CMDPacket)

	payloadReader := bytes.NewBuffer(cmdPacket.Payload)

	cid := packets.DecodeUint64(payloadReader)
	password := packets.DecodeString(payloadReader)

	client, err := m.Storage().GetClient(cid)
	if err != nil {
		m.Error("查询客户端[%d]失败！ -> %v", cid, err)
		replyCMDPacketError(m, CMDRegisterAck, RegisterError)
		return
	}
	if client != nil {
		m.Error("客户端[%d]已存在！", cid)
		replyCMDPacketError(m, CMDRegisterAck, RegisterClientExist)
		return
	}
	err = m.Storage().AddClient(tgo.NewClient(cid, password))
	if err != nil {
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
