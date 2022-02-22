package nsqconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	form "devcode.xeemore.com/systech/gojunkyard/form"
	deduplicator "devcode.xeemore.com/systech/gojunkyard/nsq/consumer/internal/deduplicator"
	panicrecover "devcode.xeemore.com/systech/gojunkyard/nsq/consumer/internal/panic"
	storage "devcode.xeemore.com/systech/gojunkyard/nsq/consumer/storage"
	nop_storage "devcode.xeemore.com/systech/gojunkyard/nsq/consumer/storage/nop"
	reporter "devcode.xeemore.com/systech/gojunkyard/reporter"
	nop_reporter "devcode.xeemore.com/systech/gojunkyard/reporter/nop"

	nsq "github.com/nsqio/go-nsq"
)

// Handler is the interface for each consumer handler
// Name must be the channel name
// Handle is the function which return "func (*type) error" used to handle the message from NSQ server
type Handler interface {
	Name() string
}

// Consumer is the which handle all of handler and config
type Consumer struct {
	config    *Config
	storage   storage.Storage
	reporter  reporter.Reporter
	handlers  map[string]Handler
	consumers []*nsq.Consumer
}

// NewConsumer will create *Consumer object
func NewConsumer(cfg *Config) *Consumer {
	return &Consumer{
		config:    cfg,
		storage:   nop_storage.New(),
		reporter:  nop_reporter.NewNopReporter(),
		handlers:  make(map[string]Handler, len(cfg.Consumers)),
		consumers: make([]*nsq.Consumer, 0, len(cfg.Consumers)),
	}
}

// SetReporter is used for report all data based on level
func (c *Consumer) SetReporter(reporter reporter.Reporter) {
	c.reporter = reporter
}

// SetStorage ...
func (c *Consumer) SetStorage(storage storage.Storage) {
	c.storage = storage
}

// RegisterHandler will add all handler which will be used at config
func (c *Consumer) RegisterHandler(h Handler) {
	if _, ok := c.handlers[h.Name()]; ok {
		panic(fmt.Sprintf("[NSQ] Handler.Handle: %s has been registered", h.Name()))
	}

	typ := reflect.ValueOf(h).MethodByName("Handle").Type()
	if typ.Kind() != reflect.Func ||
		typ.NumIn() != 3 ||
		typ.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() ||
		typ.In(1) != reflect.TypeOf((*map[string]interface{})(nil)).Elem() ||
		typ.In(2).Kind() != reflect.Ptr ||
		typ.NumOut() != 2 ||
		typ.Out(0).Kind() != reflect.Bool ||
		typ.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic(`[NSQ] Handler.Handle must be "func (ctx context.Context, tags map[string]interface{}, in *type) (bool, error)`)
	}

	c.handlers[h.Name()] = h
}

// init is used for initialize all handler based on config
func (c *Consumer) init() {
	var (
		panicrecover = panicrecover.New(c.reporter)
		deduplicator = deduplicator.New(c.reporter, c.storage)
		nsqConfig    = nsq.NewConfig()
	)

	for _, v := range c.config.Consumers {
		h, ok := c.handlers[v.Name]
		if !ok {
			c.reporter.Infof("[NSQ] Handler.Handle: %s is ignored due to not exist.\n", v.Name)
			continue
		}

		var (
			val      = reflect.ValueOf(h).MethodByName("Handle")
			typ      = val.Type()
			elem     = typ.In(2).Elem()
			isStruct = elem.Kind() == reflect.Struct
			reporter = c.reporter
		)

		for _, t := range v.Topics {
			// Don't remove this declaration!
			var (
				topic          = t
				channel        = v.Name
				skipValidation = v.SkipValidation
			)

			var h nsq.Handler = nsq.HandlerFunc(func(m *nsq.Message) error {
				var in = reflect.New(elem).Interface()

				// step 1. get request payload
				err := json.Unmarshal(m.Body, in)
				if err != nil {
					reporter.Warningf(
						"[NSQ] Consumer failed unmarshaling data. topic: %s, channel: %s, message: %s, err: %s",
						topic, channel, m.Body, err,
					)
					return nil
				}

				// step 2. do validation if it is not skipped
				if !skipValidation && isStruct {
					err = form.Validate(in)
					if err != nil {
						reporter.Warningf(
							"[NSQ] Consumer detects invalid body. topic: %s, channel: %s, message: %s, err: %s",
							topic, channel, m.Body, err,
						)
						return nil
					}
				}

				// step 3. call the value and get the (requeue and error)
				var (
					ret = val.Call([]reflect.Value{
						reflect.ValueOf(context.Background()),
						reflect.ValueOf(topic.Tags),
						reflect.ValueOf(in),
					})
					requeue = ret[0].Bool()
				)

				// step 4. if not requeue, then return err
				err, _ = ret[1].Interface().(error)
				if !requeue {
					if err != nil {
						reporter.Warningf(
							"[NSQ] Consumer detects error, but does not requeue. topic: %s, channel: %s, message: %s, err: %s",
							topic, channel, m.Body, err,
						)
						return nil
					}
					reporter.Infof(
						"[NSQ] Consumer successfully process the message. topic: %s, channel: %s, message: %s",
						topic, channel, m.Body,
					)
					return nil
				}

				// step 5. if requeue and stil have requeue attempt
				if m.Attempts <= nsqConfig.MaxAttempts {
					m.Requeue(time.Second)
					reporter.Errorf(
						"[NSQ] Consumer is requeuing the message. topic: %s, channel: %s, message: %s, err: %s",
						topic, channel, m.Body, err,
					)
					return err
				}

				// step 6. if requeue attempt is more than max attempt
				// >>> PUBLISH_MESSAGE_TO_EXCEPTION_HERE <<< //
				reporter.Errorf(
					"[NSQ] Consumer cannot requeue the message due to reaching max attempts. topic: %s, channel: %s, message: %s, err: %s",
					topic, channel, m.Body, err,
				)
				return err
			})

			consumer, err := nsq.NewConsumer(topic.Name, channel, nsqConfig)
			if err != nil {
				panic(fmt.Sprintf("[NSQ] Failed to init consumer. topic: %s, channel: %s, err: %s\n", topic, channel, err))
			}

			consumer.SetLogger(nil, nsq.LogLevelError)
			consumer.ChangeMaxInFlight(v.MaxInFlight)
			consumer.AddConcurrentHandlers(
				panicrecover(deduplicator.Handle(topic.Name, channel, h)),
				v.Concurrency,
			)

			c.consumers = append(c.consumers, consumer)
		}
	}
}

// Run will start the nsq server
func (c *Consumer) Run() error {
	c.init()
	for _, v := range c.consumers {
		err := v.ConnectToNSQLookupds(c.config.NSQLookupd)
		if err != nil {
			return err
		}
		err = v.ConnectToNSQDs(c.config.NSQD)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunGraceful is blocking version of nsq run.
// This function also doing graceful shutdown when there is SIGHUP, SIGINT, SIGTERM, SIGQUIT
func (c *Consumer) RunGraceful() error {
	defer c.Stop()

	err := c.Run()
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan

	return nil
}

// Stop is used for gracefully stop the nsq consumer
func (c *Consumer) Stop() {
	for _, v := range c.consumers {
		v.Stop()
	}
	for _, v := range c.consumers {
		<-v.StopChan
	}
}
