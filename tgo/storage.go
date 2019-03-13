package tgo

type Storage interface {
	SaveMsg(msgContext *MsgContext) error // 保存消息
	ReadMsgChan() chan *MsgContext // 读取消息
}