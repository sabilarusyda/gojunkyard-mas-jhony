package cli

import (
	"net/http"
	"os"

	"devcode.xeemore.com/systech/gojunkyard/reporter"
	"devcode.xeemore.com/systech/gojunkyard/reporter/writer"
)

type CliReporter struct {
	stdout reporter.Reporter
	stderr reporter.Reporter
}

// NewCliReporter is used to initiate slack reporter
func NewCliReporter(appName string, level writer.LEVEL) *CliReporter {
	return &CliReporter{
		stdout: writer.NewWriterReporter(appName, level, os.Stdout),
		stderr: writer.NewWriterReporter(appName, level, os.Stderr),
	}
}

func (cr *CliReporter) SetCallDepth(calldepth int) {
	type scd interface{ SetCallDepth(int) }
	if v, ok := cr.stdout.(scd); ok {
		v.SetCallDepth(calldepth)
	}
	if v, ok := cr.stderr.(scd); ok {
		v.SetCallDepth(calldepth)
	}
}

func (cr *CliReporter) SetFlags(flag int) {
	type sf interface{ SetFlags(int) }
	if v, ok := cr.stdout.(sf); ok {
		v.SetFlags(flag)
	}
	if v, ok := cr.stderr.(sf); ok {
		v.SetFlags(flag)
	}
}

func (cr *CliReporter) Debug(v ...interface{}) {
	cr.stdout.Debug(v...)
}

func (cr *CliReporter) Debugf(format string, v ...interface{}) {
	cr.stdout.Debugf(format, v...)
}

func (cr *CliReporter) Debugln(v ...interface{}) {
	cr.stdout.Debugln(v...)
}

func (cr *CliReporter) Info(v ...interface{}) {
	cr.stdout.Info(v...)
}

func (cr *CliReporter) Infof(format string, v ...interface{}) {
	cr.stdout.Infof(format, v...)
}

func (cr *CliReporter) Infoln(v ...interface{}) {
	cr.stdout.Infoln(v...)
}

func (cr *CliReporter) Warning(v ...interface{}) {
	cr.stdout.Warning(v...)
}

func (cr *CliReporter) Warningf(format string, v ...interface{}) {
	cr.stdout.Warningf(format, v...)
}

func (cr *CliReporter) Warningln(v ...interface{}) {
	cr.stdout.Warningln(v...)
}

func (cr *CliReporter) Error(v ...interface{}) {
	cr.stderr.Error(v...)
}

func (cr *CliReporter) Errorf(format string, v ...interface{}) {
	cr.stderr.Errorf(format, v...)
}

func (cr *CliReporter) Errorln(v ...interface{}) {
	cr.stderr.Errorln(v...)
}

func (cr *CliReporter) ReportPanic(err interface{}, stacktrace []byte) error {
	return cr.stderr.ReportPanic(err, stacktrace)
}

func (cr *CliReporter) ReportHTTPPanic(err interface{}, stacktrace []byte, r *http.Request) error {
	return cr.stderr.ReportHTTPPanic(err, stacktrace, r)
}
