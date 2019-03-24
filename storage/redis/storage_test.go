package redis

import (
	"github.com/tgo-team/tgo-core/test"
	"github.com/tgo-team/tgo-core/tgo"
	"testing"
)

func TestStorage_SaveAndGetMsg(t *testing.T) {
	s := NewStorage(nil)
	go func() {
		msg := <-s.StorageMsgChan()
		test.Equal(t,uint64(100),msg.ChannelID())
		test.Equal(t,uint64(1),msg.Msg().MessageID)
		test.Equal(t,uint64(2),msg.Msg().From)
	}()
	msg := tgo.NewMsg(1,2,[]byte("hello"))
	err := s.SaveMsg(tgo.NewMsgContext(msg,100))
	test.Nil(t,err)

	resultMsg,err := s.GetMsg(1)
	test.Equal(t,msg.String(),resultMsg.String())
}

func TestStorage_SaveAndGetChannel(t *testing.T) {
	s := NewStorage(nil)

	ch := &tgo.Channel{
		ChannelID:100,
		ChannelType: 1,
	}
	err := s.SaveChannel(ch)
	test.Nil(t,err)

	resultCh,err := s.GetChannel(100)
	test.Nil(t,err)

	test.Equal(t,ch.String(),resultCh.String())

}

func TestStorage_AddConsumer(t *testing.T) {

}
