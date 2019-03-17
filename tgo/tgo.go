package tgo

import (
	"sync/atomic"
)

type TGO struct {
	Servers []Server
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
	tg.Servers = GetServers(ctx)
	if tg.Servers == nil {
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
	for _,server :=range t.Servers {
		err := server.Start()
		if err!=nil {
			return err
		}
	}
	return nil
}

func (t *TGO) Stop() error {
	close(t.exitChan)
	for _,server :=range t.Servers {
		err := server.Stop()
		if err!=nil {
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
			if msgContext!=nil {
				channel, err := t.GetChannel(msgContext.ChannelID())
				if err != nil {
					t.Error("获取管道[%d]失败！-> %v",msgContext.ChannelID(),err)
					continue
				}
				if channel==nil {
					t.Error("管道[%d]不存在！",msgContext.ChannelID())
					continue
				}
				t.waitGroup.Add(1)
				go func(msgContext *MsgContext) {
					channel.DeliveryMsg(msgContext)
					t.waitGroup.Done()
				}(msgContext)
			}
		case msgContext := <-t.memoryMsgChan:
			if msgContext!=nil {
				channel, err := t.GetChannel(msgContext.ChannelID())
				if err != nil {
					t.Error("获取管道[%d]失败！-> %v",msgContext.ChannelID(),err)
					continue
				}
				if channel==nil {
					t.Error("管道[%d]不存在！",msgContext.ChannelID())
					continue
				}

				t.waitGroup.Add(1)
				go func(msgContext *MsgContext) {
					channel.DeliveryMsg(msgContext)
					t.waitGroup.Done()
				}(msgContext)
			}


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


