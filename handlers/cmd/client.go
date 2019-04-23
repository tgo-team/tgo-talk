package cmd

import (
	"bytes"
	"fmt"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
)

func UpdateToken(m *tgo.MContext) {
	m.Info("UpdateToken")
	cmdPacket := m.Packet().(*packets.CmdPacket)

	payloadReader := bytes.NewBuffer(cmdPacket.Payload)

	clientID := packets.DecodeUint64(payloadReader)
	password := packets.DecodeString(payloadReader)

	client, err := m.Storage().GetClient(clientID)
	if err != nil {
		m.Error("查询客户端[%d]失败！ -> %v", clientID, err)
		replyCmdPacketError(m, fmt.Sprintf("%d",CMDUpdateClientAck), UpdateClientError)
		return
	}
	if client == nil {
		err = m.Storage().AddClient(tgo.NewClient(clientID, password))
		if err != nil {
			replyCmdPacketError(m, fmt.Sprintf("%d",CMDUpdateClientAck), UpdateClientError)
			return
		}
	}else{
		err = m.Storage().UpdateClient(clientID, password)
		if err != nil {
			replyCmdPacketError(m, fmt.Sprintf("%d",CMDUpdateClientAck), UpdateClientError)
			return
		}
	}

	var channelID = clientID
	// 添加个人管道
	err = m.Storage().AddChannel(tgo.NewChannelModel(channelID,tgo.ChannelTypePerson))
	if err!=nil {
		m.Error("添加Channel失败！-> %v",err)
		replyCmdPacketError(m, fmt.Sprintf("%d",CMDUpdateClientAck), UpdateClientError)
		return
	}

	if err := m.Storage().Bind(clientID,channelID);err!=nil {
		m.Error("绑定Channel失败！-> %v",err)
		replyCmdPacketError(m, fmt.Sprintf("%d",CMDUpdateClientAck), UpdateClientError)
		return
	}
	replyCmdPacketSuccess(m, fmt.Sprintf("%d",CMDUpdateClientAck))

}

func replyCmdPacketError(m *tgo.MContext, cmd string, status uint16) {
	m.ReplyPacket(packets.NewCmdPacket(cmd, packets.EncodeUint16(status)))
}

func replyCmdPacketSuccess(m *tgo.MContext, cmd string) {
	m.ReplyPacket(packets.NewCmdPacket(cmd, packets.EncodeUint16(uint16(CMDSuccess))))
}
