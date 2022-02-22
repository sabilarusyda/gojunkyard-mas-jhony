package writer

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Writer struct {
	logger zerolog.Logger
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

// NewWriterReporter is used to initiate slack reporter
func NewWriterReporter(appName string, level LEVEL, w io.Writer) *Writer {
	// appName is not used anymore. for backward compatibility
	logger := zerolog.New(w).Level(zerolog.Level(level - 1))
	logger = logger.With().Timestamp().Logger()
	return &Writer{logger: logger}
}

func (w *Writer) SetCallDepth(depth int) {
	w.logger = w.logger.With().CallerWithSkipFrameCount(depth).Logger()
}

func (w *Writer) SetFlags(flag int) {
	// noop. backward compatibility
}

func (w *Writer) Debug(v ...interface{}) {
	w.logger.Debug().Msg(fmt.Sprint(v...))
}

func (w *Writer) Debugf(format string, v ...interface{}) {
	w.logger.Debug().Msg(fmt.Sprintf(format, v...))
}

func (w *Writer) Debugln(v ...interface{}) {
	w.logger.Debug().Msg(fmt.Sprint(v...))
}

func (w *Writer) Info(v ...interface{}) {
	w.logger.Info().Msg(fmt.Sprint(v...))
}

func (w *Writer) Infof(format string, v ...interface{}) {
	w.logger.Info().Msg(fmt.Sprintf(format, v...))
}

func (w *Writer) Infoln(v ...interface{}) {
	w.logger.Info().Msg(fmt.Sprint(v...))
}

func (w *Writer) Warning(v ...interface{}) {
	w.logger.Warn().Msg(fmt.Sprint(v...))
}

func (w *Writer) Warningf(format string, v ...interface{}) {
	w.logger.Warn().Msg(fmt.Sprintf(format, v...))
}

func (w *Writer) Warningln(v ...interface{}) {
	w.logger.Warn().Msg(fmt.Sprint(v...))
}

func (w *Writer) Error(v ...interface{}) {
	w.logger.Error().Msg(fmt.Sprint(v...))
}

func (w *Writer) Errorf(format string, v ...interface{}) {
	w.logger.Error().Msg(fmt.Sprintf(format, v...))
}

func (w *Writer) Errorln(v ...interface{}) {
	w.logger.Error().Msg(fmt.Sprint(v...))
}

func (w *Writer) ReportPanic(err interface{}, stacktrace []byte) error {
	w.logger.Error().Msg(fmt.Sprint(err, string(stacktrace)))
	return nil
}

func (w *Writer) ReportHTTPPanic(err interface{}, stacktrace []byte, _ *http.Request) error {
	return w.ReportPanic(err, stacktrace)
}
