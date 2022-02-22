package logger

import (
	"fmt"

	gclog "cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

// StackdriverOptions is options to setting up Stackdriver logger.
type StackdriverOptions struct {
	Client       *gclog.Client
	LoggerName   string
	LevelEnabler zapcore.LevelEnabler
}

type stackdriverCore struct {
	logger *gclog.Logger
	labels map[string]string

	zapcore.LevelEnabler
}

func newStackdriverCore(options *StackdriverOptions) zapcore.Core {
	return &stackdriverCore{
		logger:       options.Client.Logger(options.LoggerName),
		labels:       make(map[string]string),
		LevelEnabler: options.LevelEnabler,
	}
}

func (s *stackdriverCore) With(fields []zapcore.Field) zapcore.Core {
	return s.with(fields)
}

func (s *stackdriverCore) with(fields []zapcore.Field) *stackdriverCore {
	clone := s.clone()

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for i := range fields {
		fields[i].AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		clone.labels[k] = fmt.Sprint(v)
	}

	return clone
}

func (s *stackdriverCore) clone() *stackdriverCore {
	copy := *s
	return &copy
}

func (s *stackdriverCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if s.Enabled(ent.Level) {
		return ce.AddCore(ent, s)
	}
	return ce
}

func (s *stackdriverCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	clone := s.with(fields)

	entry := gclog.Entry{
		Severity: stackdriverSeverity(ent.Level),
		Labels:   clone.labels,
		Payload:  ent.Message,
	}

	go s.logger.Log(entry)

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		s.Sync()
	}

	return nil
}

func (s *stackdriverCore) Sync() error {
	return s.logger.Flush()
}

func stackdriverSeverity(lvl zapcore.Level) gclog.Severity {
	switch lvl {
	case zapcore.DebugLevel:
		return gclog.Debug
	case zapcore.InfoLevel:
		return gclog.Info
	case zapcore.WarnLevel:
		return gclog.Warning
	case zapcore.ErrorLevel:
		return gclog.Error
	case zapcore.DPanicLevel:
		return gclog.Critical
	case zapcore.PanicLevel:
		return gclog.Alert
	case zapcore.FatalLevel:
		return gclog.Emergency
	default:
		// Unrecognized levels are fatal.
		return gclog.Emergency
	}
}
