package tgo

import "sync/atomic"

type TGO struct {
	Server Server
	opts   atomic.Value // options
	*Route
	exitChan  chan int
	waitGroup WaitGroupWrapper
	Storage   Storage // storage msg
	monitor   Monitor // Monitor
}

func New(opts *Options) *TGO {
	tg := &TGO{
		exitChan: make(chan int, 0),
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

	tg.waitGroup.Wrap(tg.msgLoop)
	tg.waitGroup.Wrap(tg.storageLoop)
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
		case msg := <-t.Server.ReadMsgChan():
			if msg != nil {
				err := t.Storage.SaveMsg(msg)
				if err != nil {
					t.GetOpts().Log.Error("Failed to store message！- %v", err)
					continue
				}
			} else {
				t.GetOpts().Log.Warn("Get the message is nil")
			}
		case <-t.exitChan:
			goto exit

		}
	}
exit:
	t.GetOpts().Log.Info("msgLoop is exit!")
}

func (t *TGO) storageLoop() {
	for {
		select {
		case msg := <-t.Storage.ReceiveMsgChan():
			if msg != nil {
				t.GetOpts().Log.Info("storage the message - %v", msg)
				t.Serve(GetMContext(msg))
			} else {
				t.GetOpts().Log.Warn("storage the message is nil")
			}
		case <-t.exitChan:
			goto exit
		}
	}
exit:
	t.GetOpts().Log.Info("storageLoop is exit!")
}

func (t *TGO) TraceMsg(tag string, msgId int64) {
	t.GetOpts().Log.Info("trace [%d] is %s",msgId,tag)
}
