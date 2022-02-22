package nsqproducer

import (
	"encoding/json"
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

// AsyncProducer will handle all of producer and config asynchronous
type AsyncProducer struct {
	Producer
	transaction chan *nsq.ProducerTransaction
}

// Init is used for initialize producer
func (p *AsyncProducer) Init() {
	var nsqConfig = nsq.NewConfig()
	np, err := nsq.NewProducer(p.config.NSQD, nsqConfig)
	if err != nil {
		panic(fmt.Sprintf("[NSQ] Failed to init producer async. nsqd: %s, err: %v\n", p.config.NSQD, err))
	}

	p.producer = np
}

// MultiPublish will publish multiple data to certain topic in nsq asynchronously
func (p *AsyncProducer) MultiPublish(topic string, data []interface{}) error {
	bb := make([][]byte, 0, len(data))
	for _, v := range data {
		b, err := json.Marshal(v)
		if err != nil {
			p.reporter.Errorf(
				"[NSQ] Producer async failed marshaling data. topic: %s, message: %+v, err: %v",
				topic, v, err,
			)
			return err
		}
		bb = append(bb, b)
	}

	err := p.producer.MultiPublishAsync(topic, bb, p.transaction)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer async failed publish data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	return nil
}

// Publish will publish data to certain topic in nsq asynchronously
func (p *AsyncProducer) Publish(topic string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer async failed marshaling data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	err = p.producer.PublishAsync(topic, b, p.transaction)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer async failed publish data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	return nil
}
