package tgo

import "time"

type Options struct {
	LogLevel             LogLevel
	Log                  Log
	Monitor              Monitor
	LogPrefix            string
	Verbose              bool
	TCPAddress           string
	UDPAddress           string
	HTTPAddress          string
	HTTPSAddress         string
	MaxHeartbeatInterval time.Duration
	DataPath             string
	MaxMsgSize           int32
	MaxBytesPerFile      int64         // 每个文件数据文件最多保存多大的数据 单位byte
	SyncEvery            int64         // 内存队列每满多少消息就同步一次
	SyncTimeout          time.Duration // 超过超时时间没同步就持久化一次
	Pro                  Protocol      // 协议
	MemQueueSize         int64         // 内存队列的chan大小，值表示内存中能堆积多少条消息
	MsgTimeout           time.Duration // 消息发送超时时间
}

func NewOptions() *Options {

	return &Options{
		MaxBytesPerFile:      100 * 1024 * 1024,
		MsgTimeout:           60 * time.Second,
		MaxMsgSize:           1024 * 1024,
		MemQueueSize:         10000,
		SyncEvery:            2500,
		SyncTimeout:          2 * time.Second,
		LogPrefix:            "[tgo-server] ",
		LogLevel:             DebugLevel,
		TCPAddress:           "0.0.0.0:6666",
		UDPAddress:           "0.0.0.0:5555",
		HTTPAddress:          "0.0.0.0:6667",
		HTTPSAddress:         "0.0.0.0:6443",
		MaxHeartbeatInterval: 60 * time.Second,
		Pro:                  NewProtocol("mqtt-im"),
	}
}
