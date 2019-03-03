package log

import (
	"github.com/sirupsen/logrus"
	"github.com/tgo-team/tgo-chat/tgo"
)



type Log struct {
	logLevel tgo.LogLevel
}

func init()  {
	tgo.RegistryLog(func(logLevel tgo.LogLevel) tgo.Log {

		return NewLog(logLevel)
	})

}

func NewLog(logLevel tgo.LogLevel) *Log {
	logrus.SetLevel(logrus.Level(logLevel))
	//logrus.SetReportCaller(true)
	//logrus.SetOutput(os.Stdout)
	return &Log{
		logLevel:logLevel,
	}
}

func (lg *Log) Info(format string,a ...interface{})  {
	logrus.Infof(format,a...)
}

func (lg *Log) Error(format string,a ...interface{})  {
	logrus.Errorf(format,a...)
}

func (lg *Log) Warn(format string,a ...interface{})  {
	logrus.Warnf(format,a...)
}

func (lg *Log) Debug(format string,a ...interface{})  {
	logrus.Debugf(format,a...)
}

func (lg *Log) Fatal(format string,a ...interface{})  {
	logrus.Fatalf(format,a...)
}