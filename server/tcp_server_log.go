package server

import "fmt"
// --------- log -------------
func (s *TCPServer) Info(format string, a ...interface{}) {
	s.GetOpts().Log.Info(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *TCPServer) Error(format string, a ...interface{}) {
	s.GetOpts().Log.Error(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *TCPServer) Warn(format string, a ...interface{}) {
	s.GetOpts().Log.Warn(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *TCPServer) Debug(format string, a ...interface{}) {
	s.GetOpts().Log.Debug(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *TCPServer) Fatal(format string, a ...interface{}) {
	s.GetOpts().Log.Fatal(fmt.Sprintf("【%s】%s", s.getLogPrefix(), format), a...)
}

func (s *TCPServer) getLogPrefix() string {
	return "TCPServer"
}
