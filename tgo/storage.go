package tgo

import "github.com/tgo-team/tgo-chat/tgo/packets"

type Storage interface {
	SaveMsg(packet packets.Packet) error // 保存消息
	ReadMsgChan() chan packets.Packet // 读取消息
}