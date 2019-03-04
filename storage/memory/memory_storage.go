package memory

import (
	"fmt"
	"github.com/tgo-team/tgo-chat/tgo"
)

func init()  {

	tgo.RegistryStorage(func(context *tgo.Context) tgo.Storage {
		return NewStorage()
	})
}

type Storage struct {
	receiveMsgChan chan *tgo.Msg
	msgMap         map[string]*tgo.Msg
}

func NewStorage() *Storage {
	return &Storage{
		receiveMsgChan: make(chan *tgo.Msg, 0),
		msgMap:         make(map[string]*tgo.Msg),
	}
}

func (m *Storage) ReceiveMsgChan() chan *tgo.Msg {
	return m.receiveMsgChan
}

func (m *Storage) SaveMsg(msg *tgo.Msg) error {
	m.msgMap[fmt.Sprintf("%d-%d", msg.Id, msg.From)] = msg
	m.receiveMsgChan <- msg
	return nil
}
