package sentry

import (
	"errors"
	"fmt"
	"net/http"

	_sentry "github.com/getsentry/raven-go"
)

type Sentry struct {
	client *_sentry.Client
	option *Option
}

type Option struct {
	AppName     string  `envconfig:"APP_NAME"`
	AppVersion  string  `envconfig:"APP_VERSION"`
	DSN         string  `envconfig:"DSN"`
	Environment string  `envconfig:"ENVIRONMENT"`
	SampleRate  float32 `envconfig:"SAMPLE_RATE"`
	Level       uint8   `envconfig:"LEVEL"`
}

// NewSentryReporter is used to initiate slack reporter
func NewSentryReporter(option *Option) *Sentry {
	client, _ := _sentry.NewWithTags(option.DSN, map[string]string{
		"app_name": option.AppName,
	})
	if option.Level == 0 {
		option.Level = INFO
	}
	if option.SampleRate > 0 {
		client.SetSampleRate(option.SampleRate)
	}
	if len(option.Environment) > 0 {
		client.SetEnvironment(option.Environment)
	}
	if len(option.AppVersion) > 0 {
		client.SetRelease(option.AppVersion)
	}
	return &Sentry{
		client: client,
		option: option,
	}
}

func (s *Sentry) capture(packet *_sentry.Packet, level _sentry.Severity) {
	packet.Level = level
	s.client.Capture(packet, nil)
}

func (s *Sentry) Debug(v ...interface{}) {
	if isDebug(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprint(v...)), _sentry.DEBUG)
	}
}

func (s *Sentry) Debugf(format string, v ...interface{}) {
	if isDebug(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintf(format, v...)), _sentry.DEBUG)
	}
}

func (s *Sentry) Debugln(v ...interface{}) {
	if isDebug(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintln(v...)), _sentry.DEBUG)
	}
}

func (s *Sentry) Info(v ...interface{}) {
	if isInfo(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprint(v...)), _sentry.INFO)
	}
}

func (s *Sentry) Infof(format string, v ...interface{}) {
	if isInfo(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintf(format, v...)), _sentry.INFO)
	}
}

func (s *Sentry) Infoln(v ...interface{}) {
	if isInfo(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintln(v...)), _sentry.INFO)
	}
}

func (s *Sentry) Warning(v ...interface{}) {
	if isWarning(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprint(v...)), _sentry.WARNING)
	}
}

func (s *Sentry) Warningf(format string, v ...interface{}) {
	if isWarning(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintf(format, v...)), _sentry.WARNING)
	}
}

func (s *Sentry) Warningln(v ...interface{}) {
	if isWarning(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintln(v...)), _sentry.WARNING)
	}
}

func (s *Sentry) Error(v ...interface{}) {
	if isError(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprint(v...)), _sentry.ERROR)
	}
}

func (s *Sentry) Errorf(format string, v ...interface{}) {
	if isError(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintf(format, v...)), _sentry.ERROR)
	}
}

func (s *Sentry) Errorln(v ...interface{}) {
	if isError(s.option.Level) {
		s.capture(_sentry.NewPacket(fmt.Sprintln(v...)), _sentry.ERROR)
	}
}

// ReportPanic is used to send the panic message to slack
func (s *Sentry) ReportPanic(err interface{}, _ []byte) error {
	if isFatal(s.option.Level) {
		packet := s.generatePanicPacket(err)
		if packet != nil {
			s.client.Capture(packet, nil)
		}
	}
	return nil
}

// ReportHTTPPanic is used to send the panic message to slack
func (s *Sentry) ReportHTTPPanic(err interface{}, _ []byte, r *http.Request) error {
	if r == nil {
		return s.ReportPanic(err, nil)
	}
	if isFatal(s.option.Level) {
		packet := s.generatePanicPacket(err)
		if packet != nil {
			packet.Interfaces = append(packet.Interfaces, _sentry.NewHttp(r))
			s.client.Capture(packet, nil)
		}
	}
	return nil
}

func (s *Sentry) generatePanicPacket(err interface{}) *_sentry.Packet {
	switch rval := err.(type) {
	case nil:
		return nil
	case error:
		packet := _sentry.NewPacket(rval.Error(), _sentry.NewException(rval, _sentry.NewStacktrace(3, 3, nil)))
		packet.Level = _sentry.FATAL
		return packet
	default:
		str := fmt.Sprint(rval)
		packet := _sentry.NewPacket(str, _sentry.NewException(errors.New(str), _sentry.NewStacktrace(3, 3, nil)))
		packet.Level = _sentry.FATAL
		return packet
	}
}
