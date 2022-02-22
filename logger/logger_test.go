package logger

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	gclog "cloud.google.com/go/logging"
	sentry "github.com/getsentry/raven-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logFile = "test.log"

func TestLogger(t *testing.T) {
	ctx := context.Background()

	sdClient, err := gclog.NewClient(ctx, "staging-199507")
	if err != nil {
		return
	}

	options := &Options{
		LogStdout: true,
		Filepath:  logFile,

		Stackdriver: &StackdriverOptions{
			Client:       sdClient,
			LoggerName:   "devcode.xeemore.com/systech/gojunkyard-log-test",
			LevelEnabler: zap.LevelEnablerFunc(levelEnablerAll),
		},
	}

	logger, err := NewLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, logger)

	defer logger.Sync()

	logger.Debug("This is debug level")
	logger.Debugf("Hi, %s", "debugf")

	logger.Info("This is info level")
	logger.Infof("Hi, %s", "infof")

	logger.Warning("This is warning level")
	logger.Warningf("Hi, %s", "warningf")

	logger.Error("This is error level")
	logger.Errorf("Hi, %s", "errorf")

	logger.Critical("This is critical level")
	logger.Criticalf("Hi, %s", "criticalf")

	// NO NEED TO TEST THIS METHOD
	// logger.Emergency("This is Emergency level")
	// logger.Emergencyf("Hi, %s", "Emergencyf")

	logger = logger.With("service", "devcode.xeemore.com/systech/gojunkyard-test")

	logger.Debug("This is debug level with `service` label")
	logger.Debugf("Hi, %s", "debugf with `service` label")

	zapLogger := logger.Zap()
	assert.NotNil(t, zapLogger)
}

func TestBaseLoggerNoOptions(t *testing.T) {
	err := InitLogger(nil)
	assert.Nil(t, err)
	assert.NotNil(t, base)

	defer Sync()

	Debug("This is debug level")
	Debugf("Hi, %s", "debugf")

	Info("This is info level")
	Infof("Hi, %s", "infof")

	Warning("This is warning level")
	Warningf("Hi, %s", "warningf")

	Error("This is error level")
	Errorf("Hi, %s", "errorf")

	Critical("This is critical level")
	Criticalf("Hi, %s", "criticalf")

	With("with_with", true)

	Debug("This is debug level with `with`")
	Debugf("Hi, %s", "debugf with `with`")

	zapLogger := Zap()
	assert.NotNil(t, zapLogger)
}

func TestAlert(t *testing.T) {
	ctx := context.Background()

	sdClient, err := gclog.NewClient(ctx, "staging-199507")
	if err != nil {
		return
	}

	options := &Options{
		LogStdout: true,

		Stackdriver: &StackdriverOptions{
			Client:       sdClient,
			LoggerName:   "devcode.xeemore.com/systech/gojunkyard-log-test",
			LevelEnabler: zap.LevelEnablerFunc(levelEnablerAll),
		},
	}

	err = InitLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, base)

	defer Sync()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "This is alert level", r)
		}
	}()

	Alert("This is alert level")
}

func TestAlertf(t *testing.T) {
	ctx := context.Background()

	sdClient, err := gclog.NewClient(ctx, "staging-199507")
	if err != nil {
		return
	}

	options := &Options{
		LogStdout: true,

		Stackdriver: &StackdriverOptions{
			Client:       sdClient,
			LoggerName:   "devcode.xeemore.com/systech/gojunkyard-log-test",
			LevelEnabler: zap.LevelEnablerFunc(levelEnablerAll),
		},
	}

	err = InitLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, base)

	defer Sync()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Hi, alertf", r)
		}
	}()

	Alertf("Hi, %s", "alertf")
}

func TestWriteFile(t *testing.T) {
	if _, err := os.Stat(logFile); !os.IsNotExist(err) {
		if err := os.Remove(logFile); err != nil {
			os.Exit(1)
		}
	}

	options := &Options{
		Filepath: logFile,
	}

	logger, err := NewLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, logger)

	defer Sync()

	msg := "Valar morghulis"
	logger.Debug(msg)

	file, err := os.Open(logFile)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	assert.Contains(t, string(content), "DEBUG")
	assert.Contains(t, string(content), msg)
}

func TestInitStackdriver(t *testing.T) {
	ctx := context.Background()

	sdClient, err := gclog.NewClient(ctx, "staging-199507")
	if err != nil {
		return
	}

	// No client
	options := &Options{
		Stackdriver: &StackdriverOptions{},
	}

	logger, err := NewLogger(options)
	assert.NotNil(t, err)
	assert.Nil(t, logger)

	// No logger name
	options = &Options{
		Stackdriver: &StackdriverOptions{
			Client: sdClient,
		},
	}

	logger, err = NewLogger(options)
	assert.NotNil(t, err)
	assert.Nil(t, logger)

	// No LevelEnabler, use default
	options = &Options{
		Stackdriver: &StackdriverOptions{
			Client:     sdClient,
			LoggerName: "devcode.xeemore.com/systech/gojunkyard-log-test",
		},
	}

	logger, err = NewLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
}

func TestInitSentry(t *testing.T) {
	sentryClient, err := sentry.New("")
	if err != nil {
		return
	}

	// No client
	options := &Options{
		Sentry: &SentryOptions{},
	}

	logger, err := NewLogger(options)
	assert.NotNil(t, err)
	assert.Nil(t, logger)

	// No LevelEnabler, use default
	options = &Options{
		Sentry: &SentryOptions{
			Client: sentryClient,
		},
	}

	logger, err = NewLogger(options)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
}
