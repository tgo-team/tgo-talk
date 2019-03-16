package packets

import "fmt"

type ConnReturnCode byte

const (
	ConnReturnCodeSuccess ConnReturnCode = iota // 0x00连接已接受
	ConnReturnCodeUnsupportProtocol  //0x01连接已拒绝，不支持的协议版本
	ConnReturnCodeUnsupportClientFlag  // 0x02连接已拒绝，不合格的客户端标识符
	ConnReturnCodeUnavailableServices // 0x03连接已拒绝，服务端不可用
	ConnReturnCodePasswordOrUnameError // 0x04连接已拒绝，无效的用户名或密码
	ConnReturnCodeUnAuth // 0x05客户端未被授权连接到此服务器
	ConnReturnCodeError // 服务器内部错误
)

type ConnackPacket struct {
	FixedHeader
	ReturnCode     ConnReturnCode
}

func NewConnackPacketWithHeader(fh FixedHeader) *ConnackPacket  {
	c := &ConnackPacket{}
	c.FixedHeader = fh
	return  c
}

func NewConnackPacket(returnCode ConnReturnCode) *ConnackPacket  {
	c := &ConnackPacket{}
	c.PacketType = Connack
	c.ReturnCode = returnCode
	return c
}


func (c *ConnackPacket) GetFixedHeader() FixedHeader  {

	return c.FixedHeader
}

func (c *ConnackPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("returncode: %d", c.ReturnCode)
	return str
}