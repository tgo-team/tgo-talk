package tgo

type Server interface {
	Start() error
	MsgChan() chan Msg
	Stop() error
}
