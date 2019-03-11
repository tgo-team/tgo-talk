package packets

import "fmt"

type PingreqPacket struct {
	FixedHeader
}

func NewPingreqPacket(fh FixedHeader) *PingreqPacket  {
	pr := &PingreqPacket{}
	pr.FixedHeader = fh
	return pr
}


func (pr *PingreqPacket) GetFixedHeader() FixedHeader  {

	return pr.FixedHeader
}

func (pr *PingreqPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}
