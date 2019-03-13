package tgo

import "time"

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

//StatefulConn 有状态连接
type StatefulConn interface {
	StartIOLoop()
	SetAuth(auth bool)
	IsAuth() bool
	SetID(id uint64)
	GetID() uint64
	SetDeadline(t time.Time) error
}
