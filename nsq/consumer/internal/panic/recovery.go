package panic

import (
	"runtime/debug"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type Reporter interface {
	ReportPanic(err interface{}, stacktrace []byte) error
}

// New is used to initiate panic recover middleware
func New(rp Reporter) func(h nsq.Handler) nsq.Handler {
	return func(h nsq.Handler) nsq.Handler {
		return nsq.HandlerFunc(func(m *nsq.Message) error {
			defer func() {
				if err := recover(); err != nil {
					if rp != nil {
						rp.ReportPanic(err, debug.Stack())
					}
					m.Requeue(time.Second)
				}
			}()
			return h.HandleMessage(m)
		})
	}
}
