package server

import (
	"github.com/tgo-team/tgo-chat/test"
	"github.com/tgo-team/tgo-chat/tgo"
	"testing"
	"time"
)

func TestTCPServer_Start(t *testing.T) {
	opts := tgo.NewOptions()
	opts.Log = test.NewLog(t)

	s := NewTCPServer(opts)
	err := s.Start()
	test.Nil(t,err)
	time.Sleep(50 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)
}
