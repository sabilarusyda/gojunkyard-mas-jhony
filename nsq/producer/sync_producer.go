package nsqproducer

import (
	"encoding/json"
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

// SyncProducer will handle all of producer and config
type SyncProducer struct {
	Producer
}

// Init is used for initialize producer
func (p *SyncProducer) Init() {
	var nsqConfig = nsq.NewConfig()
	np, err := nsq.NewProducer(p.config.NSQD, nsqConfig)
	if err != nil {
		panic(fmt.Sprintf("[NSQ] Failed to init producer. nsqd: %s, err: %v\n", p.config.NSQD, err))
	}

	p.producer = np
}

// MultiPublish will publish multiple data to certain topic in nsq
func (p *SyncProducer) MultiPublish(topic string, data []interface{}) error {
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

	err := p.producer.MultiPublish(topic, bb)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer failed publish data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	return nil
}

// Publish will publish data to certain topic in nsq
func (p *SyncProducer) Publish(topic string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer failed marshaling data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	err = p.producer.Publish(topic, b)
	if err != nil {
		p.reporter.Errorf(
			"[NSQ] Producer failed publish data. topic: %s, message: %+v, err: %v",
			topic, data, err,
		)
		return err
	}

	return nil
}
