package tgo

type Server interface {
	Start() error
	Stop() error
}

//type StatefulServer interface {
//	SetDeadline(clientId int64,t time.Time) error
//	Keepalive(clientId int64) error
//	GetClient(clientId int64) Client
//	AuthClient(clientId,newClientId int64)
//	ClientIsAuth(clientId int64) bool
//}

type StatefulServer interface {
	//AddConn(clientID uint64,conn Conn) error
}