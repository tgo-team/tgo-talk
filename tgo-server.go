package main

import (
	"github.com/judwhite/go-svc/svc"
	_ "github.com/tgo-team/tgo-talk/log"
	_ "github.com/tgo-team/tgo-talk/protocol/mqtt"
	_ "github.com/tgo-team/tgo-talk/server/tcp"
	_ "github.com/tgo-team/tgo-talk/storage/memory"
	"github.com/tgo-team/tgo-talk/tgo"
	"os"
	"path/filepath"
	"syscall"
)

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		panic(err)
	}
}

type program struct {
	t *tgo.TGO
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {

	t := tgo.New(tgo.NewOptions())
	err := t.Start()
	if err != nil {
		panic(err)
	}
	//t.Use(handlers.HandleAuth)
	//t.Use(handlers.HandleHeartbeat)
	p.t = t
	return nil
}

func (p *program) Stop() error {
	if p.t != nil {
		return p.t.Stop()

	}
	return nil
}
