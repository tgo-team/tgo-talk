package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/tgo-team/tgo-chat/tgo"
	"io"
)

func init()  {
	tgo.RegistryProtocol("tgo", func() tgo.Protocol {
		return NewTGO()
	})
}


/**
  消息协议:
  | msg_tye（1 byte） | remain length（4 byte） | msg_no（8 byte）| to (8 byte) | variable header length（1 byte） | variable header | payload

 */
type TGO struct {
	Order binary.ByteOrder
}

func NewTGO() *TGO  {

	return &TGO{
		Order: binary.LittleEndian,
	}
}

func (t *TGO) Decode(reader io.Reader) (msg *tgo.Msg,err error)  {
	msg = &tgo.Msg{}
	// msg type
	var msgType int8
	if  err = binary.Read(reader,t.Order, &msgType); err != nil {
		return
	}
	msg.MsgType = tgo.MsgType(msgType)

	if msg.MsgType == tgo.MsgTypePing { // 如果是ping到此结束
		return
	}

	// remain length
	var remainLength int32
	if err = binary.Read(reader, t.Order, &remainLength); err != nil {
		return
	}

	// 消息内容
	bodyBytes := make([]byte, remainLength)
	if _, err = io.ReadFull(reader, bodyBytes); err != nil {
		return
	}
	bodyReader := bytes.NewBuffer(bodyBytes)

	// msg_no
	var msgNo int64
	if err = binary.Read(bodyReader, t.Order, &msgNo); err != nil {
		return nil, err
	}
	msg.Id = msgNo

	// to
	var to int64
	if err = binary.Read(bodyReader, t.Order, &to); err != nil {
		return nil, err
	}
	msg.To = to

	// variable header length
	var variableHeaderLength int8
	if err = binary.Read(bodyReader, t.Order, &variableHeaderLength); err != nil {
		return nil, err
	}

	// variable header
	if variableHeaderLength>0 {
		variableHeader := make([]byte,variableHeaderLength)
		if err = binary.Read(bodyReader, t.Order, &variableHeader); err != nil {
			return nil, err
		}
		msg.VariableHeader = variableHeader
	}
	// payloadLength = remainLength - msgNoLength toLength - variableHeaderLengthLength - variableHeaderLength
	payloadLength :=  int(remainLength)  - 8 - 8  - 1- int(variableHeaderLength)
	payload := make([]byte,payloadLength)
	if _,err = io.ReadFull(bodyReader, payload); err != nil {
		return nil, err
	}
	msg.Payload = payload
	return msg, nil
}

func (t *TGO) Encode(msg *tgo.Msg) ([]byte,error)  {
	if msg.MsgType == tgo.MsgTypePong { // 如果是ping到此结束
		return []byte{byte(tgo.MsgTypePong)},nil
	}
	return nil,nil
}