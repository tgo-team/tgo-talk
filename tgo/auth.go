package tgo

type Auth interface {
	Auth(clientID uint64,password string) error
}
