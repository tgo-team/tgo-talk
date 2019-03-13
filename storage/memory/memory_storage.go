package memory

import (
	"github.com/tgo-team/tgo-talk/tgo"
)

func init()  {

	tgo.RegistryStorage(func(context *tgo.Context) tgo.Storage {
		return NewStorage()
	})
}

type Storage struct {
	readMsgChan chan *tgo.MsgContext
	msgMap         map[string]*tgo.MsgContext
}

func NewStorage() *Storage {
	return &Storage{
		readMsgChan: make(chan *tgo.MsgContext, 0),
		msgMap:         make(map[string]*tgo.MsgContext),
	}
}

func (m *Storage) ReadMsgChan() chan *tgo.MsgContext {
	return m.readMsgChan
}

func (m *Storage) SaveMsg(msgContext *tgo.MsgContext) error {
	m.readMsgChan <- msgContext
	return nil
}
