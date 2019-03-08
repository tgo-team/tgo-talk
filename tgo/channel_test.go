package tgo

import (
	"github.com/tgo-team/tgo-chat/test"
	"testing"
)

func TestPersonChannel_AddConsumer(t *testing.T) {
	c := NewPersonChannel("test",getContext(t))
	c.AddConsumer("abc",NewConsumer(123,23))
	c.AddConsumer("zzz",NewConsumer(234,45))

	test.Equal(t,2, len(c.consumers))
}
