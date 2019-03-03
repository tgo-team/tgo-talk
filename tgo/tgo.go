package tgo

import "sync/atomic"

type TGO struct {
	Log
	Server Server
	opts   atomic.Value // options
	*Route
}

func New(opts *Options) *TGO {
	tg := &TGO{
	}
	if opts.Log == nil {
		tg.Log = NewLog(opts.LogLevel)
	}
	tg.storeOpts(opts)

	ctx := &Context{
		TGO: tg,
	}
	tg.Server = NewServer(ctx)

	r := NewRoute(ctx)
	tg.Route = r
	return tg
}

func (t *TGO) Start() error {
	return t.Server.Start()
}

func (t *TGO) Stop() error {

	return t.Server.Stop()
}

func (t *TGO) storeOpts(opts *Options) {
	t.opts.Store(opts)
}

func (t *TGO) GetOpts() *Options {
	return t.opts.Load().(*Options)
}
