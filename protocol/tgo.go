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
  | msg_tye（1 byte） | remain length（4 byte） | to (8 byte) | variable header length（2 byte） | variable header | payload

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
		println(err.Error())
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

	// body
	bodyBytes := make([]byte, remainLength)
	if _, err = io.ReadFull(reader, bodyBytes); err != nil {
		return
	}
	bodyReader := bytes.NewBuffer(bodyBytes)

	// to
	var to int64
	if err = binary.Read(bodyReader, t.Order, &to); err != nil {
		return nil, err
	}
	msg.To = to

	// variable header length
	var variableHeaderLength int16
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
	// payloadLength = remainLength -  toLength - variableHeaderLengthLength - variableHeaderLength
	payloadLength :=  int64(remainLength)   - 8  - 2- int64(variableHeaderLength)
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

	buff := bytes.NewBuffer(make([]byte, 0))

	// msg_type
	if err := binary.Write(buff, t.Order, int8(msg.MsgType)); err != nil {
		return nil, err
	}

	bodyBuff := bytes.NewBuffer(make([]byte,0))
	// to
	if err := binary.Write(bodyBuff, t.Order, int64(msg.To)); err != nil {
		return nil, err
	}
	// variable header length
	if err := binary.Write(bodyBuff, t.Order, int16(len(msg.VariableHeader))); err != nil {
		return nil, err
	}
	// variable header
	if len(msg.VariableHeader)>0 {
		if err := binary.Write(bodyBuff, t.Order, msg.VariableHeader); err != nil {
			return nil, err
		}
	}
	if len(msg.Payload)>0 {
		if err := binary.Write(bodyBuff, t.Order, msg.Payload); err != nil {
			return nil, err
		}
	}
	bodyBytes := bodyBuff.Bytes()
	if err := binary.Write(buff, t.Order, int32(len(bodyBytes))); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, t.Order, bodyBytes); err != nil {
		return nil, err
	}
	return buff.Bytes(),nil
}