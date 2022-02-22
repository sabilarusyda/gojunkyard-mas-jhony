package handler

import (
	"net/http"
	"sync"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

type Handler interface {
	Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int)
	RenderError(w http.ResponseWriter, r *http.Request, err error, statusCode int)
	RenderJSON(w http.ResponseWriter, r *http.Request, v interface{}, statusCode int)
}

type handler struct {
	logger      *zap.Logger
	jsonapi     jsoniter.API
	errorPool   sync.Pool
	successPool sync.Pool
}

// mimeJSON is reusable application/json type
var mimeJSON = [...]string{"application/json"}

func New() Handler {
	var logger, _ = getZapConfig().Build()
	return &handler{
		logger:      logger,
		jsonapi:     jsoniter.ConfigFastest,
		errorPool:   sync.Pool{New: func() interface{} { return new(errorResponse) }},
		successPool: sync.Pool{New: func() interface{} { return new(successResponse) }},
	}
}

func getZapConfig() *zap.Config {
	return &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.WarnLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "ts",
			LevelKey:      "level",
			CallerKey:     "caller",
			MessageKey:    "type",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func (h *handler) log(r *http.Request, body []byte, statusCode int) {
	var (
		lvl = zapLevel(statusCode)
		ce  = h.logger.Check(lvl, "RESPONSE_LOGGER")
	)
	if ce == nil {
		return
	}
	ce.Write(
		zap.Int("statusCode", statusCode),
		zap.ByteString("body", body),
		zap.String("method", r.Method),
		zap.String("path", r.RequestURI),
	)
}

func zapLevel(statusCode int) zapcore.Level {
	switch {
	case statusCode < 300:
		return zapcore.InfoLevel
	case statusCode < 400:
		return zapcore.InfoLevel
	case statusCode < 500:
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func (h *handler) Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	// step 1. write the redirect url
	http.Redirect(w, r, url, statusCode)

	// step 2. write to the log
	h.log(r, []byte(url), statusCode)
}

func (h *handler) RenderError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	// step 1. write to header
	header := w.Header()
	header["Content-Type"] = mimeJSON[:]

	// step 2. fill the response
	response := h.errorPool.Get().(*errorResponse)
	response.Error = http.StatusText(statusCode)

	// step 3. write the response
	w.WriteHeader(statusCode)
	h.jsonapi.NewEncoder(w).Encode(response)

	// step 4. prepare for the log
	response.Error = getErrorMessage(err)
	byt, _ := h.jsonapi.Marshal(response)
	h.log(r, byt, statusCode)

	// step 5. put the error pool back
	h.errorPool.Put(response)
}

func (h *handler) RenderJSON(w http.ResponseWriter, r *http.Request, v interface{}, statusCode int) {
	// step 1. write to the header
	header := w.Header()
	header["Content-Type"] = mimeJSON[:]

	// step 2. fill the response
	response := h.successPool.Get().(*successResponse)
	response.Data = v

	// step 3. prepare for the response
	byt, _ := h.jsonapi.Marshal(response)

	// step 4. write to the response
	w.WriteHeader(statusCode)
	w.Write(byt)

	// step 5. write to the log
	h.log(r, byt, statusCode)

	// step 6. put the success pool back
	h.successPool.Put(response)
}

func getErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
