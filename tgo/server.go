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
	GetClient(clientId int64) Client
	SetClientAuthInfo(clientId int64,authId int64,token string) error
}