package tgo

import (
	"sync/atomic"
)

type TGO struct {
	Server Server
	opts   atomic.Value // options
	*Route
	exitChan        chan int
	waitGroup       WaitGroupWrapper
	Storage         Storage // storage msg
	monitor         Monitor // Monitor
	channelMap      map[uint64]*Channel
	memoryMsgChan   chan *MsgContext
	ConnContextChan chan *ConnContext
	ConnManager *connManager
}

func New(opts *Options) *TGO {
	tg := &TGO{
		exitChan:        make(chan int, 0),
		channelMap:      map[uint64]*Channel{},
		memoryMsgChan:   make(chan *MsgContext, opts.MemQueueSize),
		ConnContextChan: make(chan *ConnContext, 1024),
		ConnManager: newConnManager(),
	}
	if opts.Log == nil {
		opts.Log = NewLog(opts.LogLevel)
	}
	//if opts.Monitor == nil {
	//	opts.Monitor = tg
	//}
	tg.storeOpts(opts)

	ctx := &Context{
		TGO: tg,
	}

	// server
	tg.Server = NewServer(ctx)
	if tg.Server == nil {
		opts.Log.Fatal("请先配置Server！")
	}

	// route
	tg.Route = NewRoute(ctx)

	// storage
	tg.Storage = NewStorage(ctx) // new storage
	if tg.Storage == nil {
		opts.Log.Fatal("请先配置存储！")
	}

	tg.waitGroup.Wrap(tg.msgLoop)
	return tg
}

func (t *TGO) Start() error {
	return t.Server.Start()
}

func (t *TGO) Stop() error {
	close(t.exitChan)
	if t.Server != nil {
		err := t.Server.Stop()
		if err != nil {
			return err
		}
	}
	t.waitGroup.Wait()
	t.Info("TGO -> 退出")
	return nil
}

func (t *TGO) storeOpts(opts *Options) {
	t.opts.Store(opts)
}

func (t *TGO) GetOpts() *Options {
	return t.opts.Load().(*Options)
}

func (t *TGO) msgLoop() {
	for {
		select {
		case connContext := <-t.ConnContextChan:
			if connContext != nil {
				t.Info("收到消息 -> %v", connContext)
				t.Serve(GetMContext(connContext))
			} else {
				t.Warn("Receive the message is nil")
			}
		case msgContext := <-t.Storage.StorageMsgChan():
			if msgContext != nil {
				t.GetOpts().Log.Info("Storage-ReceiveMsgChan--%v", msgContext)
				channel, err := t.GetChannel(msgContext.ChannelID())
				if err != nil {

				}
				println(channel)
				//t.StartInFlightTimeout(msg, 0)
			}
		case msgContext := <-t.memoryMsgChan:
			t.startPushMsg(msgContext)

		case <-t.exitChan:
			goto exit

		}
	}
exit:
	t.Debug("停止收取消息。")
}

func (t *TGO) GetChannel(channelID uint64) (*Channel, error) {
	channel, ok := t.channelMap[channelID]
	var err error
	if !ok {
		channel, err = t.Storage.GetChannel(channelID)
		if err != nil {
			return nil, err
		}
		if channel != nil {
			t.channelMap[channelID] = channel
		}
	}
	return channel, nil
}

func (t *TGO) startPushMsg(msgCtx *MsgContext)  {
	t.Debug("将消息[%d]下发到管道[%d]！",msgCtx.Msg().MessageID,msgCtx.channelID)
	clientIDs,err := t.Storage.GetClientIDs(msgCtx.channelID)
	if err!=nil {
		t.Error("获取管道[%d]的客户端ID集合失败！ -> %v",err)
		return
	}
	for _,clientID :=range clientIDs {
		if clientID == msgCtx.Msg().From { // 不发送给自己
			continue
		}
		online := IsOnline(clientID)
		if online {

		}
	}

}
