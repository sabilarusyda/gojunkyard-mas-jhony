package writer

import (
	"github.com/stretchr/testify/mock"
)

type _writer struct {
	mock.Mock
}

func (w *_writer) Write(p []byte) (n int, err error) {
	args := w.Called(p)
	return args.Int(0), args.Error(1)
}
