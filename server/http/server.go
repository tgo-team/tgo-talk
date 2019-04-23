package http

import (
	"bytes"
	"fmt"
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
	case "/client/update":
		s.clientUpdate(w, req)
	}

}

func (s *Server) clientUpdate(w http.ResponseWriter, req *http.Request) {
	var clientUpdateStruct struct {
		ClientID uint64 `json:"client_id"`
		Password string `json:"password"`
	}
	if err := utils.BindJson(req, &clientUpdateStruct); err != nil {
		s.Error("数据输入有误！-> %v", err)
		utils.ResponseSuccessJson(w,map[string]interface{}{
			"status": 400,
			"msg": "数据输入有误！",
		})
		return
	}
	var body bytes.Buffer
	body.Write(packets.EncodeUint64(clientUpdateStruct.ClientID))
	body.Write(packets.EncodeString(clientUpdateStruct.Password))
	CmdPacket := packets.NewCmdPacket(fmt.Sprintf("%d",cmd.CMDUpdateClient), body.Bytes())

	respPacket,err := s.requestCMD(req,CmdPacket)
	if err!=nil {
		s.Error("请求出错！-> %v",err)
		utils.ResponseSuccessJson(w,map[string]interface{}{
			"status": 400,
			"msg": "执行命令出错！",
		})
		return
	}

	status := packets.DecodeUint16(bytes.NewBuffer(respPacket.Payload))
	if status == cmd.UpdateClientError {
		utils.ResponseSuccessJson(w,map[string]interface{}{
			"status": 400,
			"msg": "更新客户端出错！",
		})
		return
	}
	utils.ResponseSuccessJson(w,map[string]interface{}{
		"status": 200,
	})
}

func (s *Server) requestCMD(req *http.Request, packet *packets.CmdPacket) (*packets.CmdPacket,error) {
	respChan := make(chan []byte, 0)
	cn := NewConn(req, respChan, NewConnChan(s.ctx.TGO.AcceptPacketChan, nil), s.ctx)
	s.ctx.TGO.AcceptPacketChan <- tgo.NewPacketContext(packet, cn)

	var respPacket *packets.CmdPacket
	select {
	case data := <-respChan:
		packet,err := s.ctx.TGO.GetOpts().Pro.DecodePacket(bytes.NewBuffer(data))
		if err!=nil {
			return nil,err
		}
		respPacket = packet.(*packets.CmdPacket)
	case <-time.After(5 * time.Second): //超时5s
	}
	return respPacket,nil
}
