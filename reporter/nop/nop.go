package nop

import "net/http"

type Nop bool

func NewNopReporter() *Nop {
	return new(Nop)
}

func (*Nop) Debug(v ...interface{})                                                    {}
func (*Nop) Debugf(format string, v ...interface{})                                    {}
func (*Nop) Debugln(v ...interface{})                                                  {}
func (*Nop) Info(v ...interface{})                                                     {}
func (*Nop) Infof(format string, v ...interface{})                                     {}
func (*Nop) Infoln(v ...interface{})                                                   {}
func (*Nop) Warning(v ...interface{})                                                  {}
func (*Nop) Warningf(format string, v ...interface{})                                  {}
func (*Nop) Warningln(v ...interface{})                                                {}
func (*Nop) Error(v ...interface{})                                                    {}
func (*Nop) Errorf(format string, v ...interface{})                                    {}
func (*Nop) Errorln(v ...interface{})                                                  {}
func (*Nop) ReportPanic(err interface{}, stacktrace []byte) error                      { return nil }
func (*Nop) ReportHTTPPanic(err interface{}, stacktrace []byte, r *http.Request) error { return nil }
