package packets

import "fmt"

type ConnectPacket struct {
	FixedHeader
	ClientID uint64
	UsernameFlag     bool
	PasswordFlag     bool
	Username         string
	Password         string

	Keepalive       uint16
}

func NewConnectPacketWithHeader(fh FixedHeader) *ConnectPacket  {
	c := &ConnectPacket{}
	c.FixedHeader=fh
	return  c
}

func NewConnectPacket(clientID uint64,password string) *ConnectPacket   {
	c := &ConnectPacket{}
	c.PacketType = Connect
	c.ClientID = clientID
	c.Password = password
	return c
}

func (c *ConnectPacket) GetFixedHeader() FixedHeader  {

	return c.FixedHeader
}

func (c *ConnectPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("Usernameflag: %t Passwordflag: %t keepalive: %d clientId: %d Username: %s Password: %s", c.UsernameFlag, c.PasswordFlag, c.Keepalive, c.ClientID, c.Username, c.Password)
	return str
}