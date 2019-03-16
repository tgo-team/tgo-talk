package cmd

import (
	"github.com/tgo-team/tgo-talk/tgo"
	"github.com/tgo-team/tgo-talk/tgo/packets"
)

func Register(m *tgo.MContext)  {
	m.Info("Register")
	err := m.ReplyPacket(packets.NewCMDPacket(2,[]byte{200}))
	if err!=nil {
		panic(err)
	}

}