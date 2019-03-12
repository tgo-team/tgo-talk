package tgo

import "fmt"

// --------- log -------------
func (t *TGO) Info(format string, a ...interface{}) {
	t.GetOpts().Log.Info(fmt.Sprintf("【%s】%s", t.getLogPrefix(), format), a...)
}

func (t *TGO) Error(format string, a ...interface{}) {
	t.GetOpts().Log.Error(fmt.Sprintf("【%s】%s", t.getLogPrefix(), format), a...)
}

func (t *TGO) Warn(format string, a ...interface{}) {
	t.GetOpts().Log.Warn(fmt.Sprintf("【%s】%s", t.getLogPrefix(), format), a...)
}

func (t *TGO) Debug(format string, a ...interface{}) {
	t.GetOpts().Log.Debug(fmt.Sprintf("【%s】%s", t.getLogPrefix(), format), a...)
}

func (t *TGO) Fatal(format string, a ...interface{}) {
	t.GetOpts().Log.Fatal(fmt.Sprintf("【%s】%s", t.getLogPrefix(), format), a...)
}

func (t *TGO) getLogPrefix() string {
	return "TGO"
}
