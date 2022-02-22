package deduplicator

import (
	"time"

	"devcode.xeemore.com/systech/gojunkyard/nsq/consumer/storage"
	"devcode.xeemore.com/systech/gojunkyard/reporter"

	nsq "github.com/nsqio/go-nsq"
)

type Deduplicator struct {
	reporter reporter.Reporter
	storage  storage.Storage
}

func New(reporter reporter.Reporter, storage storage.Storage) *Deduplicator {
	return &Deduplicator{
		reporter: reporter,
		storage:  storage,
	}
}

func (d *Deduplicator) Handle(topic, channel string, h nsq.Handler) nsq.Handler {
	return nsq.HandlerFunc(func(m *nsq.Message) error {
		var key = calculateKey(topic, channel, m.Body)

		// 1. Set the key if not exist
		ok, err := d.storage.SetNX(key, 3*time.Minute)
		if err != nil {
			d.reporter.Errorf(
				"[NSQ_DEDUPLICATOR] Failed to set the key. topic: %s, channel: %s, message: %s, err: %s",
				topic, channel, m.Body, err,
			)
			m.Requeue(time.Second)
			return err
		}
		// 2. if key has been exist, then return
		if !ok {
			d.reporter.Warningf(
				"[NSQ_DEDUPLICATOR] Message has been processed, topic: %s, channel: %s, ignoring the message. message: %s",
				topic, channel, m.Body, err,
			)
			return nil
		}
		// 3. call the handler
		err = h.HandleMessage(m)
		if err == nil {
			return nil
		}
		// 4. delete from storage if it is error
		err = d.storage.Delete(key)
		if err != nil {
			d.reporter.Errorf(
				"[NSQ_DEDUPLICATOR] Consumer is requeue but cannot delete the message deduplicator.topic: %s, channel: %s, message: %s, err: %s",
				topic, channel, m.Body, err,
			)
		}
		return err
	})
}
