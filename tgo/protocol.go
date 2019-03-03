package tgo

import "io"

type Protocol interface {
	//Decode 解码
	Decode(reader io.Reader) (*Msg,error)
	//Encode 编码
	Encode(msg *Msg) ([]byte,error)
}