package tgo

import "time"

type Server interface {
	Start() error
	ReadMsgChan() chan *Msg
	SendMsg(to int64,msg *Msg) error
	Stop() error
}

type StatefulServer interface {
	SetDeadline(clientId int64,t time.Time) error
	Keepalive(clientId int64) error
	GetClient(clientId int64) Client
	AuthClient(clientId,newClientId int64)
	ClientIsAuth(clientId int64) bool
}