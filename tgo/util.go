package tgo

import (
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"github.com/zheng-ji/goCuckoo"
)

var cuckooFilter *cuckoo.Filter
func init()  {
	cuckooFilter = cuckoo.NewFilter(10000)
}

// 设置客户端是否上线
func Online(clientID uint64,online int) bool  {
	if online==1 {
		return cuckooFilter.Insert(packets.EncodeUint64(clientID))
	}else{
		return cuckooFilter.Del(packets.EncodeUint64(clientID))
	}
}

// 是否在线
func IsOnline(clientID uint64) bool  {
	return cuckooFilter.Find(packets.EncodeUint64(clientID))
}