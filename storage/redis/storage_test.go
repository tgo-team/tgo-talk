package redis

import (
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-talk/test"
	"testing"
)

func TestStorage_AddAndGetMsg(t *testing.T) {
	s := NewStorage(nil)
	go func() {
		msg := <-s.StorageMsgChan()
		test.Equal(t, uint64(100), msg.ChannelID())
		test.Equal(t, uint64(1), msg.Msg().MessageID)
		test.Equal(t, uint64(2), msg.Msg().From)
	}()
	msg := tgo.NewMsg(1, 2, []byte("hello"))
	err := s.AddMsg(tgo.NewMsgContext(msg, 100))
	test.Nil(t, err)

	resultMsg, err := s.GetMsg(1)
	test.Equal(t, msg.String(), resultMsg.String())
}

func TestStorage_AddAndGetChannel(t *testing.T) {
	s := NewStorage(nil)

	ch := &tgo.Channel{
		ChannelID:   100,
		ChannelType: 1,
	}
	err := s.AddChannel(ch)
	test.Nil(t, err)

	resultCh, err := s.GetChannel(100)
	test.Nil(t, err)

	test.Equal(t, ch.String(), resultCh.String())
}


func TestStorage_GetMsgWithChannel(t *testing.T) {
	s := NewStorage(nil)
	go func() {
		for {
			<-s.StorageMsgChan()
		}

	}()
	for i:=10;i<20;i++ {
		msg := tgo.NewMsg(uint64(1+i), 2, []byte("hello"))
		err := s.AddMsg(tgo.NewMsgContext(msg, 100))
		test.Nil(t,err)
	}
	msgs,err := s.GetMsgWithChannel(100,1,10)
	test.Nil(t,err)
	test.Equal(t,10,len(msgs))

}