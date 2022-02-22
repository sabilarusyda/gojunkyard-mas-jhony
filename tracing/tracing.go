package tracing

import (
	"io"

	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// Options for the tracer.
type Options struct {
	JaegerConfig jaegerconfig.Configuration
}

// Tracer wraps jaeger tracer
type Tracer struct {
	closer io.Closer
}

// NewTracer returns an initialized Tracer.
func NewTracer(name string, options *Options) (*Tracer, error) {
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := options.JaegerConfig.InitGlobalTracer(
		name,
		jaegerconfig.Metrics(jMetricsFactory),
	)
	if err != nil {
		return nil, err
	}

	tracer := &Tracer{
		closer: closer,
	}
	return tracer, nil
}

// Close closes tracer.
func (t *Tracer) Close() error {
	return t.closer.Close()
}
