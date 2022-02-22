package nsqconsumer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	nop_storage "devcode.xeemore.com/systech/gojunkyard/nsq/consumer/storage/nop"
	nop_reporter "devcode.xeemore.com/systech/gojunkyard/reporter/nop"

	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

type handlerMock struct {
	mock.Mock
}

func (hm *handlerMock) Name() string {
	return hm.Called().String(0)
}

func (hm *handlerMock) Handle(ctx context.Context, tags map[string]interface{}, in *struct {
	ID int64 `json:"id"`
}) (bool, error) {
	args := hm.Called(tags, in)
	return args.Bool(0), args.Error(1)
}

type invalidHandlerMock struct{}

func (ihm *invalidHandlerMock) Name() string {
	return ""
}

func (ihm *invalidHandlerMock) Handle() string {
	return ""
}

func TestNewConsumer(t *testing.T) {
	var (
		cfg  = NewConfig(nil, nil)
		got  = NewConsumer(cfg)
		want = &Consumer{
			config:    cfg,
			handlers:  make(map[string]Handler, len(cfg.Consumers)),
			consumers: make([]*nsq.Consumer, 0, len(cfg.Consumers)),
			reporter:  nop_reporter.NewNopReporter(),
			storage:   nop_storage.New(),
		}
	)
	assert.Equal(t, want, got)
}

func TestConsumer_SetReporter(t *testing.T) {
	consumer := NewConsumer(&Config{})
	assert.NotNil(t, consumer.reporter)

	consumer.SetReporter(nil)
	assert.Nil(t, consumer.reporter)
}

func TestConsumer_SetStorage(t *testing.T) {
	consumer := NewConsumer(&Config{})
	assert.NotNil(t, consumer.storage)

	consumer.SetStorage(nil)
	assert.Nil(t, consumer.storage)
}

func TestConsumer_RegisterHandler(t *testing.T) {
	consumer := NewConsumer(&Config{})
	assert.NotNil(t, consumer.storage)

	hm := new(handlerMock)
	hm.On("Name").Return("CONSUMER_INSERT_TO_SALESFORCE")

	// Must be panic if redeclared
	assert.Panics(t, func() { consumer.RegisterHandler(new(invalidHandlerMock)) })

	// Must not panic if first declare
	assert.NotPanics(t, func() { consumer.RegisterHandler(hm) })

	// Must be panic if redeclared
	assert.Panics(t, func() { consumer.RegisterHandler(hm) })
}
