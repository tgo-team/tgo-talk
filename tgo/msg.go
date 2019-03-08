package tgo

import (
	"fmt"
	"time"
)

type MsgType int8

const (
	MsgTypeAuth MsgType = iota
	MsgTypeAuthACK
	MsgTypeSend
	MsgTypeSendACK
	MsgTypeReceive
	MsgTypeReceiveACK
	MsgTypePing
	MsgTypePong
)

type AuthStatus int8
const (
	MsgStatusAuthFail AuthStatus = iota
	MsgStatusAuthOk
)

type MsgStatus int8
const (
	MsgStatusFail MsgStatus = iota
	MsgStatusSuccess
)

type Msg struct {
	MsgData
	ClientId int64
	index int // 在队列里的下标 内部参数
	Match string // 匹配规则
}

type MsgData struct {
	Id        int64 // Message ID
	MsgType   MsgType   // Message Type
	VariableHeader []byte // 可变头
	Payload   []byte    // Payload 消息内容
	Timestamp int64     //// 消息创建时间戳
	Attempts  uint16    // 尝试次数
	UID  int64     // 投递源ID
	ToUID   int64     // 投递目标ID
	// for in-flight handling
	DeliveryTS time.Time     // 投递超时时间（超过这个时间后就不投递了）
	Priority   int64         // 优先级 数字越小越优先
	Deferred   time.Duration // 延迟时间（延迟消息有效）

}


func NewPong() *Msg  {
	return  &Msg{
		MsgData: MsgData{
			MsgType:MsgTypePong,
		},
	}
}

func NewAuthACK(status AuthStatus) *Msg {
	return  &Msg{
		MsgData: MsgData{
			MsgType:MsgTypeAuthACK,
			VariableHeader: []byte{byte(status)},
		},
	}
}

func NewSendMsgACK(msgId int64,status MsgStatus) *Msg {
	return  &Msg{
		MsgData: MsgData{
			MsgType:MsgTypeSendACK,
			VariableHeader: []byte{byte(status)},
		},
	}
}


func (m *Msg) String() string {

	return fmt.Sprintf("Message ID:%v MsgType: %d Payload %+v", m.Id, m.MsgType, m.Payload)
}

