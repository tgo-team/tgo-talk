package packets

import "fmt"

type ConnectPacket struct {
	FixedHeader
	ClientIdentifier uint64
	UsernameFlag     bool
	PasswordFlag     bool
	Username         string
	Password         []byte

	Keepalive       uint16
}

func NewConnectPacket(fh FixedHeader) *ConnectPacket  {
	c := &ConnectPacket{}
	c.FixedHeader=fh
	return  c
}

func (c *ConnectPacket) GetFixedHeader() FixedHeader  {

	return c.FixedHeader
}

func (c *ConnectPacket) String() string {
	str := fmt.Sprintf("%s", c.FixedHeader)
	str += " "
	str += fmt.Sprintf("Usernameflag: %t Passwordflag: %t keepalive: %d clientId: %d Username: %s Password: %s", c.UsernameFlag, c.PasswordFlag, c.Keepalive, c.ClientIdentifier, c.Username, c.Password)
	return str
}