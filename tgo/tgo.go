package tgo

import "sync/atomic"

type TGO struct {
	Server Server
	opts   atomic.Value // options
	*Route
	exitChan chan int
	waitGroup    WaitGroupWrapper
}

func New(opts *Options) *TGO {
	tg := &TGO{
		exitChan :make(chan int,0),
	}
	if opts.Log == nil {
		opts.Log = NewLog(opts.LogLevel)
	}
	tg.storeOpts(opts)

	ctx := &Context{
		TGO: tg,
	}
	tg.Server = NewServer(ctx)

	r := NewRoute(ctx)
	tg.Route = r

	tg.waitGroup.Wrap(tg.msgLoop)
	return tg
}

func (t *TGO) Start() error {
	return t.Server.Start()
}

func (t *TGO) Stop() error {
	close(t.exitChan)
	if t.Server!=nil {
		err := t.Server.Stop()
		if err!=nil {
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
			if msg!=nil {
				t.GetOpts().Log.Info("Get the message - %v", msg)
				t.Serve(GetMContext(msg))
			}else{
				t.GetOpts().Log.Warn("Get the message is nil")
			}
		case <-t.exitChan:
			goto exit

		}
	}
exit:
	t.GetOpts().Log.Info("msgLoop is exit!")
}
