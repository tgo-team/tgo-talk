package tgo

import (
	"sync/atomic"
)

type TGO struct {
	Server Server
	opts   atomic.Value // options
	*Route
	exitChan         chan int
	waitGroup        WaitGroupWrapper
	Storage          Storage // storage msg
	monitor          Monitor // Monitor
	channelMap map[uint64]*Channel
	memoryMsgChan chan *MsgContext
	ConnContextChan chan *ConnContext
}

func New(opts *Options) *TGO {
	tg := &TGO{
		exitChan:         make(chan int, 0),
		channelMap: map[uint64]*Channel{},
		memoryMsgChan: make(chan *MsgContext, opts.MemQueueSize),
		ConnContextChan: make(chan *ConnContext, 1024),
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

	//tg.initPQ()

	tg.waitGroup.Wrap(tg.msgLoop)
	return tg
}

//func (t *TGO) initPQ() {
//	pqSize := int(math.Max(1, float64(t.ctx.TGO.GetOpts().MemQueueSize)/10))
//
//	t.inFlightMutex.Lock()
//	t.inFlightMessages = make(map[MsgID]*Msg)
//	t.inFlightPQ = newInFlightPqueue(pqSize)
//	t.inFlightMutex.Unlock()
//}

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
		case packet := <-t.Storage.ReadMsgChan():
			t.GetOpts().Log.Info("Storage-ReceiveMsgChan--%v", packet)
			//t.StartInFlightTimeout(msg, 0)
		case msgContext := <-t.memoryMsgChan:
			t.GetOpts().Log.Info("MemoryMsgChan--%v", msgContext)

		case <-t.exitChan:
			goto exit

		}
	}
exit:
	t.Debug("停止收取消息。")
}

func (t *TGO) GetChannel(channelID uint64) *Channel {
	channel, ok := t.channelMap[channelID]
	if !ok {
		channel = NewChannel(channelID, t.ctx)
		t.channelMap[channelID] = channel
	}
	return channel
}
