package tgo

type Monitor interface {
	Counter(flag string,inc int64)
}
