package memory

import (
	"github.com/tgo-team/tgo-core/tgo"
)

func init() {
	tgo.RegistryStorage(func(context *tgo.Context) tgo.Storage {
		return NewStorage(context)
	})
}

type Storage struct {
	storageMsgChan chan *tgo.MsgContext
	channelMsgMap  map[uint64][]*tgo.Msg
	channelMap     map[uint64]*tgo.ChannelModel
	clientMap      map[uint64]*tgo.Client
	clientChannelRelationMap  map[uint64][]uint64
	ctx            *tgo.Context
}

func NewStorage(ctx *tgo.Context) *Storage {
	return &Storage{
		storageMsgChan: make(chan *tgo.MsgContext, 0),
		channelMsgMap:  make(map[uint64][]*tgo.Msg),
		channelMap:     make(map[uint64]*tgo.ChannelModel),
		clientMap:      make(map[uint64]*tgo.Client),
		clientChannelRelationMap: make(map[uint64][]uint64),
		ctx:            ctx,
	}
}

func (s *Storage) StorageMsgChan() chan *tgo.MsgContext {
	return s.storageMsgChan
}

func (s *Storage) AddMsgInChannel(msg *tgo.Msg,channelID uint64) error {
	msgs := s.channelMsgMap[channelID]
	if msgs == nil {
		msgs = make([]*tgo.Msg, 0)
	}
	msgs = append(msgs, msg)
	s.channelMsgMap[channelID] = msgs
	s.storageMsgChan <- tgo.NewMsgContext(msg,channelID)
	return nil
}

func (s *Storage) AddChannel(c *tgo.ChannelModel) error {
	s.channelMap[c.ChannelID] = c
	return nil
}
func (s *Storage) GetChannel(channelID uint64) (*tgo.ChannelModel, error) {
	ch := s.channelMap[channelID]
	return ch, nil
}

func (s *Storage) AddClient(c *tgo.Client) error {
	s.clientMap[c.ClientID] = c
	return nil
}

func (s *Storage) Bind(clientID uint64, channelID uint64) error {
	clientIDs := s.clientChannelRelationMap[channelID]
	if clientIDs==nil {
		clientIDs = make([]uint64,0)
	}
	clientIDs = append(clientIDs,clientID)
	s.clientChannelRelationMap[channelID] = clientIDs
	return nil
}

func (s *Storage) GetClientIDs(channelID uint64) ([]uint64, error) {
	return s.clientChannelRelationMap[channelID], nil
}

func (s *Storage) GetClient(clientID uint64) (*tgo.Client, error) {

	return s.clientMap[clientID], nil
}

func (s *Storage) GetMsgInChannel(channelID uint64, pageIndex int64, pageSize int64) ([]*tgo.Msg, error) {
	msgList := s.channelMsgMap[channelID]
	if int64(len(msgList)) >= (pageIndex-1)*pageSize+pageSize {
		return msgList[(pageIndex-1)*pageSize : (pageIndex-1)*pageSize+pageSize], nil
	}
	return msgList[(pageIndex-1)*pageSize:], nil
}

func (s *Storage) RemoveMsgInChannel(messageIDs []uint64, channelID uint64)   error {

	return nil
}