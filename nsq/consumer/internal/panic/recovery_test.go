package panic

import (
	"testing"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/mock"
)

type _reporter struct {
	mock.Mock
}

func (r *_reporter) ReportPanic(err interface{}, stacktrace []byte) error {
	return r.Called(err, nil).Error(0)
}

type _delegator struct {
	mock.Mock
}

func (d *_delegator) OnFinish(*nsq.Message) {
	d.Called()
}

func (d *_delegator) OnRequeue(*nsq.Message, time.Duration, bool) {
	d.Called()
}

func (d *_delegator) OnTouch(*nsq.Message) {
	d.Called()
}

func TestRecoverWithReporter(t *testing.T) {
	// step 1. prepare testing variable
	var (
		reporter  = new(_reporter)
		delegator = new(_delegator)
		mw        = New(reporter)
		f         = nsq.HandlerFunc(func(m *nsq.Message) error {
			panic("PANIC!!!")
		})
	)

	// step 2. prepare mocking object (we dont need to test stacktrace because is will change for every execution)
	reporter.On("ReportPanic", "PANIC!!!", nil).Return(nil)
	delegator.On("OnRequeue").Once()

	// step 3. call the function wanted to test
	mw(f).HandleMessage(&nsq.Message{
		Delegate: delegator,
	})
}

func TestRecoverWithoutReporter(t *testing.T) {
	// step 1. prepare testing variable
	var (
		mw        = New(nil)
		delegator = new(_delegator)
		f         = nsq.HandlerFunc(func(m *nsq.Message) error {
			panic("PANIC")
		})
	)

	// step 2. prepare mocking object (we dont need to test stacktrace because is will change for every execution)
	delegator.On("OnRequeue").Once()

	// step 3. call the function wanted to test
	mw(f).HandleMessage(&nsq.Message{
		Delegate: delegator,
	})
}
