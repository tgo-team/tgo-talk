package tgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

// --------- message -------------

type Msg struct {
	MessageID uint64 // 消息唯一编号
	From      uint64 // 发送者ID
	Payload   []byte // 消息内容
}

func NewMsg(messageID uint64,from uint64, payload []byte) *Msg {

	return &Msg{
		MessageID: messageID,
		From: from,
		Payload:   payload,
	}
}

func (m *Msg) String() string {

	return fmt.Sprintf("MessageID: %d From: %d Payload: %s",m.MessageID,m.From,string(m.Payload))
}

func (m *Msg) MarshalBinary() (data []byte, err error) {
	var body bytes.Buffer
	body.Write(packets.EncodeUint64(m.From))
	body.Write(packets.EncodeUint64(m.MessageID))
	body.Write(m.Payload)
	return body.Bytes(),nil
}

func (m *Msg) UnmarshalBinary(data []byte) error {
	m.From = binary.BigEndian.Uint64(data[:8])
	m.MessageID = binary.BigEndian.Uint64(data[8:16])
	m.Payload = data[16:]
	return nil
}



// -------- MsgContext ------------
type MsgContext struct {
	msg *Msg
	channelID uint64
}

func NewMsgContext(msg *Msg,channelID uint64) *MsgContext {

	return &MsgContext{msg:msg,channelID:channelID}
}

func (mc *MsgContext) Msg() *Msg {
	return mc.msg
}

func (mc *MsgContext) ChannelID() uint64 {
	return mc.channelID
}

