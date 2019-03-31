package http

import "fmt"
// --------- log -------------
func (s *Server) Info(format string, a ...interface{}) {
	s.ctx.TGO.GetOpts().Log.Info(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Error(format string, a ...interface{}) {
	s.ctx.TGO.GetOpts().Log.Error(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Warn(format string, a ...interface{}) {
	s.ctx.TGO.GetOpts().Log.Warn(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Debug(format string, a ...interface{}) {
	s.ctx.TGO.GetOpts().Log.Debug(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) Fatal(format string, a ...interface{}) {
	s.ctx.TGO.GetOpts().Log.Fatal(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *Server) getLogPrefix() string {
	return "HTTPServer"
}
