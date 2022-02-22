package reporter

import "net/http"

type Reporter interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Warningln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	ReportPanic(err interface{}, stacktrace []byte) error
	ReportHTTPPanic(err interface{}, stacktrace []byte, r *http.Request) error
}
