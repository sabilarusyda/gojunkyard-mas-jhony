package aggregator

import (
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/reporter"
)

type Aggregator struct {
	rs []reporter.Reporter
}

func NewAggregator(rs ...reporter.Reporter) *Aggregator {
	return &Aggregator{rs}
}

func (a *Aggregator) Debug(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Debug(v...)
	}
}

func (a *Aggregator) Debugf(format string, v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Debugf(format, v...)
	}
}

func (a *Aggregator) Debugln(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Debugln(v...)
	}
}

func (a *Aggregator) Info(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Info(v...)
	}
}

func (a *Aggregator) Infof(format string, v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Infof(format, v...)
	}
}

func (a *Aggregator) Infoln(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Infoln(v...)
	}
}

func (a *Aggregator) Warning(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Warning(v...)
	}
}

func (a *Aggregator) Warningf(format string, v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Warningf(format, v...)
	}
}

func (a *Aggregator) Warningln(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Warningln(v...)
	}
}

func (a *Aggregator) Error(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Error(v...)
	}
}

func (a *Aggregator) Errorf(format string, v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Errorf(format, v...)
	}
}

func (a *Aggregator) Errorln(v ...interface{}) {
	for _, reporter := range a.rs {
		reporter.Errorln(v...)
	}
}

func (a *Aggregator) ReportPanic(err interface{}, stacktrace []byte) error {
	for _, reporter := range a.rs {
		reporter.ReportPanic(err, stacktrace)
	}
	return nil
}

func (a *Aggregator) ReportHTTPPanic(err interface{}, stacktrace []byte, r *http.Request) error {
	for _, reporter := range a.rs {
		reporter.ReportHTTPPanic(err, stacktrace, r)
	}
	return nil
}
