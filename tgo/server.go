package tgo

type Server interface {
	Start() error
	ReadMsgChan() chan *Msg
	WriteMsgChan() chan *Msg
	Stop() error
}
