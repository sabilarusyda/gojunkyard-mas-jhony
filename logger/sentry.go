package logger

import (
	sentry "github.com/getsentry/raven-go"
	"go.uber.org/zap/zapcore"
)

// SentryOptions is options to setting up Sentry logger.
type SentryOptions struct {
	Client       *sentry.Client
	Tags         map[string]string
	WithTrace    bool
	LevelEnabler zapcore.LevelEnabler
}

type sentryCore struct {
	client    *sentry.Client
	fields    map[string]interface{}
	tags      map[string]string
	withTrace bool

	zapcore.LevelEnabler
}

func newSentryCore(options *SentryOptions) zapcore.Core {
	return &sentryCore{
		client:       options.Client,
		fields:       make(map[string]interface{}),
		tags:         options.Tags,
		withTrace:    options.WithTrace,
		LevelEnabler: options.LevelEnabler,
	}
}

func (s *sentryCore) With(fields []zapcore.Field) zapcore.Core {
	return s.with(fields)
}

func (s *sentryCore) with(fields []zapcore.Field) *sentryCore {
	// Copy our map.
	m := make(map[string]interface{}, len(s.fields))
	for k, v := range s.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for i := range fields {
		fields[i].AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return &sentryCore{
		client:       s.client,
		LevelEnabler: s.LevelEnabler,
	}
}

func (s *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if s.Enabled(ent.Level) {
		return ce.AddCore(ent, s)
	}
	return ce
}

func (s *sentryCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	clone := s.with(fields)

	packet := &sentry.Packet{
		Message:     ent.Message,
		Timestamp:   sentry.Timestamp(ent.Time),
		Level:       sentrySeverity(ent.Level),
		Platform:    "go",
		Extra:       clone.fields,
		Fingerprint: []string{ent.Message},
	}

	if s.withTrace {
		trace := sentry.NewStacktrace(2, 3, nil)
		if trace != nil {
			packet.Interfaces = append(packet.Interfaces, trace)
		}
	}

	go clone.client.Capture(packet, s.tags)

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		clone.Sync()
	}

	return nil
}

func (s *sentryCore) Sync() error {
	s.client.Wait()
	return nil
}

func sentrySeverity(lvl zapcore.Level) sentry.Severity {
	switch lvl {
	case zapcore.DebugLevel:
		return sentry.INFO
	case zapcore.InfoLevel:
		return sentry.INFO
	case zapcore.WarnLevel:
		return sentry.WARNING
	case zapcore.ErrorLevel:
		return sentry.ERROR
	case zapcore.DPanicLevel:
		return sentry.FATAL
	case zapcore.PanicLevel:
		return sentry.FATAL
	case zapcore.FatalLevel:
		return sentry.FATAL
	default:
		// Unrecognized levels are fatal.
		return sentry.FATAL
	}
}
