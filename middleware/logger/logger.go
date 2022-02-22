package lm

import (
	"net"
	"net/http"
	"os"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/middleware/requestid"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

// New ...
func New(level uint8) func(http.Handler) http.Handler {
	logger := zerolog.New(os.Stdout).Level(zerolog.Level(level - 1))
	logger = logger.With().Timestamp().Logger()

	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		var (
			remote     = getRemoteAddr(r)
			statusCode = getStatusCode(status)
			event      = getEvent(&logger, statusCode)
		)

		if id := requestid.GetFromRequest(r); len(id) > 0 {
			event.Str("request_id", id)
		}

		event.
			Str("remote", remote).
			Str("method", r.Method).
			Str("path", r.URL.String()).
			Int("code", statusCode).
			Int("size", size).
			Dur("duration", duration)

		if origin := r.Header.Get("origin"); len(origin) > 0 {
			event.Str("origin", origin)
		}

		if ua := r.Header.Get("user-agent"); len(ua) > 0 {
			event.Str("agent", ua)
		}

		event.Send()
	})
}

func getStatusCode(status int) int {
	if status == 0 {
		return http.StatusOK
	}
	return status
}

func getRemoteAddr(r *http.Request) string {
	ip := r.Header.Get("x-forwarded-for")
	if len(ip) > 0 {
		return ip
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getEvent(logger *zerolog.Logger, status int) *zerolog.Event {
	switch {
	case status < http.StatusBadRequest:
		return logger.Info()
	case status < http.StatusInternalServerError:
		return logger.Warn()
	default:
		return logger.Error()
	}
}
