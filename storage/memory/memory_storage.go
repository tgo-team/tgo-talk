package memory

import (
	"github.com/tgo-team/tgo-chat/tgo"
	"github.com/tgo-team/tgo-chat/tgo/packets"
)

func init()  {

	tgo.RegistryStorage(func(context *tgo.Context) tgo.Storage {
		return NewStorage()
	})
}

type Storage struct {
	readMsgChan chan packets.Packet
	msgMap         map[string]packets.Packet
}

func NewStorage() *Storage {
	return &Storage{
		readMsgChan: make(chan packets.Packet, 0),
		msgMap:         make(map[string]packets.Packet),
	}
}

func (m *Storage) ReadMsgChan() chan packets.Packet {
	return m.readMsgChan
}

func (m *Storage) SaveMsg(msg packets.Packet) error {
	m.readMsgChan <- msg
	return nil
}
