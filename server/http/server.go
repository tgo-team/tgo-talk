package http

import (
	"bytes"
	"github.com/tgo-team/tgo-core/tgo"
	"github.com/tgo-team/tgo-core/tgo/packets"
	"github.com/tgo-team/tgo-talk/handlers/cmd"
	"github.com/tgo-team/tgo-talk/utils"
	"net/http"
	"time"
)

func init() {
	tgo.RegistryServer(func(context *tgo.Context) tgo.Server {
		return NewServer(context)

	})
}

type Server struct {
	ctx *tgo.Context
}

func NewServer(ctx *tgo.Context) *Server {
	return &Server{ctx: ctx}
}

func (s *Server) Start() error {
	s.Info("开始监听 -> %s", s.ctx.TGO.GetOpts().HTTPAddress)
	go func() {
		http.Handle("/", s)
		err := http.ListenAndServe(s.ctx.TGO.GetOpts().HTTPAddress, nil)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	case "/register":
		s.register(w, req)
	}

}

func (s *Server) register(w http.ResponseWriter, req *http.Request) {
	var registerStruct struct {
		ClientID uint64 `json:"client_id"`
		Password string `json:"password"`
	}
	if err := utils.BindJson(req, &registerStruct); err != nil {
		s.Error("数据输入有误！-> %v", err)
		utils.ResponseError400(w, "数据输入有误！")
		return
	}
	var body bytes.Buffer
	body.Write(packets.EncodeUint64(registerStruct.ClientID))
	body.Write(packets.EncodeString(registerStruct.Password))
	cmdPacket := packets.NewCMDPacket(cmd.CMDRegister, body.Bytes())

	respPacket,err := s.requestCMD(req,cmdPacket)
	if err!=nil {
		s.Error("请求出错！-> %v",err)
		utils.ResponseError400(w,"请求出错！")
		return
	}

	status := packets.DecodeUint16(bytes.NewBuffer(respPacket.Payload))
	if status == cmd.RegisterClientExist {
		utils.ResponseError400(w,"客户端已存在！")
		return
	}
	if status == cmd.RegisterError {
		utils.ResponseError400(w,"注册出错！")
		return
	}
	utils.ResponseSuccess(w)
}

func (s *Server) requestCMD(req *http.Request, packet *packets.CMDPacket) (*packets.CMDPacket,error) {
	respChan := make(chan []byte, 0)
	cn := NewConn(req, respChan, NewConnChan(s.ctx.TGO.AcceptPacketChan, nil), s.ctx)
	s.ctx.TGO.AcceptPacketChan <- tgo.NewPacketContext(packet, cn)

	var respPacket *packets.CMDPacket
	select {
	case data := <-respChan:
		packet,err := s.ctx.TGO.GetOpts().Pro.DecodePacket(bytes.NewBuffer(data))
		if err!=nil {
			return nil,err
		}
		respPacket = packet.(*packets.CMDPacket)
	case <-time.After(5 * time.Second): //超时5s
	}
	return respPacket,nil
}
