package tgo

import (
	"github.com/tgo-team/tgo-chat/tgo/packets"
	"time"
)

type Server interface {
	Start() error
	ReceivePacketChan() chan packets.Packet
	SendMsg(to uint64,packet packets.Packet) error
	Stop() error
}

type StatefulServer interface {
	SetDeadline(clientId int64,t time.Time) error
	Keepalive(clientId int64) error
	GetClient(clientId int64) Client
	AuthClient(clientId,newClientId int64)
	ClientIsAuth(clientId int64) bool
}