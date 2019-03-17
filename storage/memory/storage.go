package memory

import (
	"github.com/tgo-team/tgo-talk/tgo"
)

func init() {
	tgo.RegistryStorage(func(context *tgo.Context) tgo.Storage {
		return NewStorage(context)
	})
}

type Storage struct {
	storageMsgChan chan *tgo.MsgContext
	channelMsgMap      map[uint64][]*tgo.Msg
	channelMap map[uint64] *tgo.Channel
	clientMap map[uint64] *tgo.Client
	ctx *tgo.Context
}

func NewStorage(ctx *tgo.Context) *Storage {
	return &Storage{
		storageMsgChan: make(chan *tgo.MsgContext, 0),
		channelMsgMap:      make(map[uint64][]*tgo.Msg),
		channelMap: make(map[uint64]*tgo.Channel),
		clientMap: make(map[uint64]*tgo.Client),
		ctx: ctx,
	}
}

func (s *Storage) StorageMsgChan() chan *tgo.MsgContext {
	return s.storageMsgChan
}

func (s *Storage) AddMsg(msgContext *tgo.MsgContext) error {
	msgs := s.channelMsgMap[msgContext.ChannelID()]
	if msgs==nil  {
		msgs = make([]*tgo.Msg,0)
	}
	msgs = append(msgs,msgContext.Msg())
	s.channelMsgMap[msgContext.ChannelID()] = msgs
	s.storageMsgChan <- msgContext
	return nil
}

func (s *Storage) AddChannel(c *tgo.Channel) error {
	s.channelMap[c.ChannelID] = c
	return nil
}
func (s *Storage) GetChannel(channelID uint64) (*tgo.Channel,error) {
	ch := s.channelMap[channelID]
	ch.Ctx = s.ctx
	return ch,nil
}

func (s *Storage) AddClient( c *tgo.Client) error {
	s.clientMap[c.ClientID] = c
	return nil
}

func (s *Storage) Bind(consumerID uint64, channelID uint64) error {
	return nil
}

func (s *Storage) GetClientIDs(channelID uint64) ([]uint64 ,error) {
	return nil,nil
}

func (s *Storage) GetClient(clientID uint64) (*tgo.Client,error){

	return s.clientMap[clientID],nil
}