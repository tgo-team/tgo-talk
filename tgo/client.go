package tgo

type Client interface {
	Write(b []byte) error
	Exit() error
}
