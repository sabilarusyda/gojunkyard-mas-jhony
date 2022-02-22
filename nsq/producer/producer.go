package nsqproducer

import (
	"devcode.xeemore.com/systech/gojunkyard/reporter"
	"devcode.xeemore.com/systech/gojunkyard/reporter/nop"

	nsq "github.com/nsqio/go-nsq"
)

// IProducer interface of producer
type IProducer interface {
	Init()
	Publish(topic string, data interface{}) error
	MultiPublish(topic string, data []interface{}) error
}

// Producer will handle all of producer and config
type Producer struct {
	config   *Config
	reporter reporter.Reporter
	producer *nsq.Producer
}

// NewProducer will create Producer object
func NewProducer(cfg *Config) IProducer {
	if cfg.IsAsync {
		return &AsyncProducer{
			Producer: Producer{
				config:   cfg,
				reporter: nop.NewNopReporter(),
				producer: new(nsq.Producer),
			},
			transaction: make(chan *nsq.ProducerTransaction, 100),
		}
	}

	return &SyncProducer{
		Producer: Producer{
			config:   cfg,
			reporter: nop.NewNopReporter(),
			producer: new(nsq.Producer),
		},
	}
}
