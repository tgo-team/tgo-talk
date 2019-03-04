package tgo

type Storage interface {
	SaveMsg(msg *Msg) error // 保存消息
	ReceiveMsgChan() chan *Msg // 读取消息
}