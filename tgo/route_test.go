package tgo



type ServerTest struct {

}


func (s *ServerTest) Start() error {
	return nil
}
func (s *ServerTest)  ReceiveMsgChan() chan *Msg {
	return nil
}
func (s *ServerTest)  SendMsg(to int64,msg *Msg) error {
	return nil
}
func (s *ServerTest)  Stop() error {
	return nil
}

type StorageTest struct {

}

func (s *StorageTest) SaveMsg(msg *Msg) error {
	return nil
}
func (s *StorageTest) ReceiveMsgChan() chan *Msg {
	return nil
}