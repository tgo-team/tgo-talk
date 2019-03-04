package test

import (
	"testing"
)



type Log struct {
	t testing.TB
}


func NewLog(t testing.TB) *Log {
	//logrus.SetReportCaller(true)
	//logrus.SetOutput(os.Stdout)
	return &Log{
		t:t,
	}
}

func (lg *Log) Info(format string,a ...interface{})  {
	lg.t.Logf(format,a...)
}

func (lg *Log) Error(format string,a ...interface{})  {
	lg.t.Errorf(format,a...)
}

func (lg *Log) Warn(format string,a ...interface{})  {
}

func (lg *Log) Debug(format string,a ...interface{})  {
}

func (lg *Log) Fatal(format string,a ...interface{})  {
	lg.t.Fatalf(format,a...)
}