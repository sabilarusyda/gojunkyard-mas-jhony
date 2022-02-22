package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options for the logger.
type Options struct {
	LogStdout bool
	Filepath  string

	Stackdriver *StackdriverOptions
	Sentry      *SentryOptions
}

// Logger is the interface for logger used in the application components.
type Logger interface {
	With(key string, value interface{}) Logger

	Zap() *zap.Logger
	Sync() error

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})

	Alert(args ...interface{})
	Alertf(format string, args ...interface{})

	Emergency(args ...interface{})
	Emergencyf(format string, args ...interface{})
}

// logger object.
type logger struct {
	z *zap.Logger
}

func levelEnablerHigh(lvl zapcore.Level) bool {
	return lvl >= zapcore.ErrorLevel
}

func levelEnablerLow(lvl zapcore.Level) bool {
	return lvl < zapcore.ErrorLevel
}

func levelEnablerAll(lvl zapcore.Level) bool {
	return true
}

// NewLogger returns initalized Logger.
func NewLogger(options *Options) (Logger, error) {
	jsonEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	var zapCores []zapcore.Core

	if options != nil {
		if options.LogStdout {
			debugWriter := zapcore.Lock(os.Stdout)
			zapCores = append(zapCores, zapcore.NewCore(jsonEncoder, debugWriter, zap.LevelEnablerFunc(levelEnablerLow)))

			errorWriter := zapcore.Lock(os.Stderr)
			zapCores = append(zapCores, zapcore.NewCore(jsonEncoder, errorWriter, zap.LevelEnablerFunc(levelEnablerHigh)))
		}

		if options.Filepath != "" {
			if err := os.MkdirAll(filepath.Dir(options.Filepath), 0700); err != nil {
				return nil, fmt.Errorf("Failed to create log file (%s)", options.Filepath)
			}

			if _, err := os.OpenFile(options.Filepath, os.O_CREATE, 0666); err != nil {
				return nil, fmt.Errorf("Failed to open log file (%s)", options.Filepath)
			}

			writer, _, err := zap.Open(options.Filepath)
			if err != nil {
				return nil, fmt.Errorf("Failed to create writer to file (%s)", options.Filepath)
			}

			zapCores = append(zapCores, zapcore.NewCore(jsonEncoder, writer, zap.LevelEnablerFunc(levelEnablerAll)))
		}

		if options.Stackdriver != nil {
			if options.Stackdriver.Client == nil {
				return nil, fmt.Errorf("Stackdriver client is nil")
			}

			if options.Stackdriver.LoggerName == "" {
				return nil, fmt.Errorf("Stackdriver logger name is missing")
			}

			if options.Stackdriver.LevelEnabler == nil {
				options.Stackdriver.LevelEnabler = zap.LevelEnablerFunc(levelEnablerHigh)
			}

			sdOptions := &StackdriverOptions{
				Client:       options.Stackdriver.Client,
				LoggerName:   options.Stackdriver.LoggerName,
				LevelEnabler: options.Stackdriver.LevelEnabler,
			}
			sdCore := newStackdriverCore(sdOptions)
			zapCores = append(zapCores, sdCore)
		}

		if options.Sentry != nil {
			if options.Sentry.Client == nil {
				return nil, fmt.Errorf("Sentry client is nil")
			}

			if options.Sentry.LevelEnabler == nil {
				options.Sentry.LevelEnabler = zap.LevelEnablerFunc(levelEnablerHigh)
			}

			sentryOptions := &SentryOptions{
				Client:       options.Sentry.Client,
				Tags:         options.Sentry.Tags,
				WithTrace:    options.Sentry.WithTrace,
				LevelEnabler: options.Sentry.LevelEnabler,
			}
			sentryCore := newSentryCore(sentryOptions)
			zapCores = append(zapCores, sentryCore)
		}
	}

	zapTee := zapcore.NewTee(zapCores...)
	z := zap.New(zapTee)

	logger := logger{
		z: z,
	}
	return logger, nil
}

// With adds key and value to log.
func (l logger) With(key string, value interface{}) Logger {
	return logger{
		z: l.z.With(zap.String(key, fmt.Sprint(value))),
	}
}

// Zap returns *zap.Logger.
func (l logger) Zap() *zap.Logger {
	return l.z
}

// Sync flushes the buffered logs.
func (l logger) Sync() error {
	return l.z.Sync()
}

// Debug logs a message at level Debug.
func (l logger) Debug(args ...interface{}) {
	l.z.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug.
func (l logger) Debugf(format string, args ...interface{}) {
	l.z.Debug(fmt.Sprintf(format, args...))
}

//Info logs a message at level Info.
func (l logger) Info(args ...interface{}) {
	l.z.Info(fmt.Sprint(args...))
}

//Infof logs a message at level Info.
func (l logger) Infof(format string, args ...interface{}) {
	l.z.Info(fmt.Sprintf(format, args...))
}

//Warning logs a message at level Warning.
func (l logger) Warning(args ...interface{}) {
	l.z.Warn(fmt.Sprint(args...))
}

//Warningf logs a message at level Warning.
func (l logger) Warningf(format string, args ...interface{}) {
	l.z.Warn(fmt.Sprintf(format, args...))
}

//Error logs a message at level Error.
func (l logger) Error(args ...interface{}) {
	l.z.Error(fmt.Sprint(args...))
}

//Errorf logs a message at level Error.
func (l logger) Errorf(format string, args ...interface{}) {
	l.z.Error(fmt.Sprintf(format, args...))
}

//Critical logs a message at level Critical.
func (l logger) Critical(args ...interface{}) {
	l.z.DPanic(fmt.Sprint(args...))
}

//Criticalf logs a message at level Critical.
func (l logger) Criticalf(format string, args ...interface{}) {
	l.z.DPanic(fmt.Sprintf(format, args...))
}

//Alert logs a message at level Alert.
func (l logger) Alert(args ...interface{}) {
	l.z.Panic(fmt.Sprint(args...))
}

//Alertf logs a message at level Alert.
func (l logger) Alertf(format string, args ...interface{}) {
	l.z.Panic(fmt.Sprintf(format, args...))
}

//Emergency logs a message at level Emergency.
func (l logger) Emergency(args ...interface{}) {
	l.z.Fatal(fmt.Sprint(args...))
}

//Emergencyf logs a message at level Emergency.
func (l logger) Emergencyf(format string, args ...interface{}) {
	l.z.Fatal(fmt.Sprintf(format, args...))
}

var base logger

// InitLogger initializes global logger.
func InitLogger(options *Options) error {
	b, err := NewLogger(options)
	if b != nil {
		base = b.(logger)
	}
	return err
}

// With adds key and value to log.
func With(key string, value interface{}) {
	base = base.With(key, value).(logger)
}

// Zap returns *zap.Logger.
func Zap() *zap.Logger {
	return base.z
}

// Sync flushes the buffered logs.
func Sync() error {
	return base.Sync()
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	base.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug.
func Debugf(format string, args ...interface{}) {
	base.Debug(fmt.Sprintf(format, args...))
}

//Info logs a message at level Info.
func Info(args ...interface{}) {
	base.Info(fmt.Sprint(args...))
}

//Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	base.Info(fmt.Sprintf(format, args...))
}

//Warning logs a message at level Warning.
func Warning(args ...interface{}) {
	base.Warning(fmt.Sprint(args...))
}

//Warningf logs a message at level Warning.
func Warningf(format string, args ...interface{}) {
	base.Warningf(fmt.Sprintf(format, args...))
}

//Error logs a message at level Error.
func Error(args ...interface{}) {
	base.Error(fmt.Sprint(args...))
}

//Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	base.Errorf(fmt.Sprintf(format, args...))
}

//Critical logs a message at level Critical.
func Critical(args ...interface{}) {
	base.Critical(fmt.Sprint(args...))
}

//Criticalf logs a message at level Critical.
func Criticalf(format string, args ...interface{}) {
	base.Criticalf(fmt.Sprintf(format, args...))
}

//Alert logs a message at level Alert.
func Alert(args ...interface{}) {
	base.Alert(fmt.Sprint(args...))
}

//Alertf logs a message at level Alert.
func Alertf(format string, args ...interface{}) {
	base.Alertf(fmt.Sprintf(format, args...))
}

//Emergency logs a message at level Emergency.
func Emergency(args ...interface{}) {
	base.Emergency(fmt.Sprint(args...))
}

//Emergencyf logs a message at level Emergency.
func Emergencyf(format string, args ...interface{}) {
	base.Emergencyf(fmt.Sprintf(format, args...))
}
