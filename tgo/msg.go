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

type Msg struct {
	MsgData
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
	From  int64     // 投递源ID
	To    int64     // 投递目标ID
	// for in-flight handling
	DeliveryTS time.Time     // 投递超时时间（超过这个时间后就不投递了）
	Priority   int64         // 优先级 数字越小越优先
	Deferred   time.Duration // 延迟时间（延迟消息有效）

}

func (m *Msg) String() string {

	return fmt.Sprintf("Message ID:%v MsgType: %d Payload %+v", m.Id, m.MsgType, m.Payload)
}