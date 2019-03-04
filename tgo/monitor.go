package tgo

type Monitor interface {
	TraceMsg(tag string,msgId int64)
}
