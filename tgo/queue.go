package tgo

type SendQueue struct {

}

func NewSendQueue(channelID uint64,ctx *Context) *SendQueue {

	return &SendQueue{}
}

func Send(consumerIDs []uint64,msg *Msg)  {

}
