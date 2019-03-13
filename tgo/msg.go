package tgo

// --------- message -------------

type Msg struct {
	From      uint64 // 发送者ID
	MessageID uint64 // 消息唯一编号
	Payload   []byte // 消息内容
}

func NewMsg(messageID uint64,from uint64, payload []byte) *Msg {

	return &Msg{
		MessageID: messageID,
		From: from,
		Payload:   payload,
	}
}

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