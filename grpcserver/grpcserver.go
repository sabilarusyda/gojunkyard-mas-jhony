package grpcserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

type GRPC struct {
	grpc    *grpc.Server
	options *Options
}

type Options struct {
	ListenAddress         string        `envconfig:"LISTEN_ADDRESS"`
	MaxConnectionIdle     time.Duration `envconfig:"MAX_CONNECTION_IDLE"`
	MaxConnectionAge      time.Duration `envconfig:"MAX_CONNECTION_AGE"`
	MaxConnectionAgeGrace time.Duration `envconfig:"MAX_CONNECTION_AGE_GRACE"`
	MaxKeepaliveAge       time.Duration `envconfig:"MAX_KEEPALIVE_AGE"`
	MinKeepaliveAge       time.Duration `envconfig:"MIN_KEEPALIVE_AGE"`
	Timeout               time.Duration `envconfig:"TIMEOUT"`
	PermitWithoutStream   bool          `envconfig:"PERMIT_WITHOUT_STREAM"`
}

func New(options *Options, grpcOpts ...grpc.ServerOption) *GRPC {
	return &GRPC{
		options: options,
		grpc:    grpc.NewServer(grpcOpts...),
	}
}

// Run serves the HTTP endpoints.
func (g *GRPC) Run() chan error {
	var ch = make(chan error)
	go g.run(ch)
	return ch
}

func (g *GRPC) run(ch chan error) {
	listener, err := net.Listen("tcp", g.options.ListenAddress)
	if err != nil {
		ch <- err
		return
	}

	ch <- g.grpc.Serve(listener)
}

func (g *GRPC) Server() *grpc.Server {
	return g.grpc
}

func (g *GRPC) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	g.grpc.RegisterService(sd, ss)
}

// Stop terminate the server gracefully
func (g *GRPC) Stop() {
	g.grpc.GracefulStop()
}
