package tgo

import (
	"github.com/tgo-team/tgo-chat/tgo/packets"
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
	groupChannelMap  map[string]Channel
	personChannelMap map[string]Channel
	memoryPacketChan    chan packets.Packet

	Auth Auth
}

func New(opts *Options) *TGO {
	tg := &TGO{
		exitChan:         make(chan int, 0),
		groupChannelMap:  map[string]Channel{},
		personChannelMap: map[string]Channel{},
		memoryPacketChan:    make(chan packets.Packet, opts.MemQueueSize),
	}
	if opts.Log == nil {
		opts.Log = NewLog(opts.LogLevel)
	}
	if opts.Monitor == nil {
		opts.Monitor = tg
	}
	tg.storeOpts(opts)

	ctx := &Context{
		TGO: tg,
	}
	tg.Server = NewServer(ctx)
	if tg.Server == nil {
		opts.Log.Fatal("You have not configured server yet!")
	}

	tg.Route = NewRoute(ctx)     // new route
	tg.Storage = NewStorage(ctx) // new storage
	if tg.Storage == nil {
		opts.Log.Fatal("You have not configured storage yet!")
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
		case packet := <-t.Server.ReceivePacketChan():
			if packet != nil {
				t.GetOpts().Log.Info("Receive the message - %v", packet)
				t.Serve(GetMContext(packet))
			} else {
				t.GetOpts().Log.Warn("Receive the message is nil")
			}
		case packet := <-t.Storage.ReadMsgChan():
			t.GetOpts().Log.Info("Storage-ReceiveMsgChan--%v", packet)
			//t.StartInFlightTimeout(msg, 0)
		case packet := <-t.memoryPacketChan:
			t.GetOpts().Log.Info("MemoryMsgChan--%v", packet)

		case <-t.exitChan:
			goto exit

		}
	}
exit:
	t.GetOpts().Log.Info("msgLoop is exit!")
}


func (t *TGO) GetChannel(channelName string, channelType ChannelType) Channel {
	var channel Channel
	var ok bool
	if channelType == ChannelTypePerson {
		channel, ok = t.personChannelMap[channelName]
		if !ok {
			channel = NewPersonChannel(channelName, t.ctx)
			t.personChannelMap[channelName] = channel
		}
	} else if channelType == ChannelTypeGroup {
		channel, ok = t.groupChannelMap[channelName]
		if !ok {
			channel = NewGroupChannel(channelName, t.ctx)
			t.groupChannelMap[channelName] = channel
		}
	}
	return channel
}

