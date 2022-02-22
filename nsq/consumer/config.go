package nsqconsumer

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

// Config holds all configuration of nsqd, lookupd, and consumer
type Config struct {
	NSQD       []string        `envconfig:"NSQD"`
	NSQLookupd []string        `envconfig:"NSQLOOKUPD"`
	Consumers  consumerConfigs `envconfig:"CONSUMERS"`
}

// NewConfig returns pointer of Config object
func NewConfig(daddr, lookupdaddr []string) *Config {
	return &Config{
		NSQD:       daddr,
		NSQLookupd: lookupdaddr,
		Consumers:  make(consumerConfigs, 0),
	}
}

// AddConsumerConfig append the configuration of consumer using pointer of ConsumerConfig struct
func (c *Config) AddConsumerConfig(cc *ConsumerConfig) {
	if cc != nil {
		c.Consumers = append(c.Consumers, cc)
	}
}

// ConsumerConfig holds all configuration about consumer
// Name: channel name
// Topic: topic that wanted to subscribe
// MaxInFlight: number of message that will be pulled from server for every call
// Concurrency: number of asyncronous handler that will process the message from server
// SkipValidation: if you does not want to validate the data, by toggle this field to true will increase the performance
type ConsumerConfig struct {
	Name           string  `json:"name"`
	Topics         []Topic `json:"topics"`
	MaxInFlight    int     `json:"maxInFlight"`
	Concurrency    int     `json:"concurrency"`
	SkipValidation bool    `json:"skipValidation"`
	Deduplicator   bool    `json:"deduplicator"`
}

// NewConsumerConfig is factory that is used to create ConsumerConfig object.
func NewConsumerConfig(
	name string,
	topics []Topic,
	maxInFlight int,
	concurrency int,
	skipValidation bool,
	deduplicator bool,
) *ConsumerConfig {
	return &ConsumerConfig{
		Name:           name,
		Topics:         topics,
		MaxInFlight:    maxInFlight,
		Concurrency:    concurrency,
		SkipValidation: skipValidation,
		Deduplicator:   deduplicator,
	}
}

// consumerConfigs is slice of ConsumerConfig pointer. This alias is used to implement envconfig
type consumerConfigs []*ConsumerConfig

// Decode is one of envconfig interface implementation for decoding env variables
func (cc *consumerConfigs) Decode(value string) error {
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = (*reflect.SliceHeader)(unsafe.Pointer(&value)).Data
	bh.Len = len(value)
	bh.Cap = len(value)
	return json.Unmarshal(b, cc)
}

// Topic is struct that contains name of topic would be consumed and tags sent to the consumer.
// Sometime we use same handle to subscribe to some topics, but there is some special handler to multiplex the request
// The tags can be used to handle that condition
// Example:
// * We have consumer that registered to 2 topic, but we use same consumer handler. to identify where the message come from, we can use tag
// ** Topic 1. add tag {"source": "topic-1", "expireDuration": "1h"}
// ** Topic 2. add tag {"source": "topic-2", "expireDuration": "2h"}
type Topic struct {
	Name string
	Tags map[string]interface{}
}

// NewTopic ...
func NewTopic(name string, tags map[string]interface{}) *Topic {
	return &Topic{
		Name: name,
		Tags: tags,
	}
}

// UnmarshalJSON is used to handle backward compatibility configuration
// Acceptable topic parameter
// 1. "ACCOUNTS_DEVICE_REGISTRATION"
// 2. {"name": "ACCOUNTS_DEVICE_REGISTRATION", "tags": {"source": "accounts-device-registration"}}
func (t *Topic) UnmarshalJSON(b []byte) error {
	var v struct {
		Name string                 `json:"name"`
		Tags map[string]interface{} `json:"tags"`
	}

	err := json.Unmarshal(b, &v)
	if err == nil {
		t.Name = v.Name
		t.Tags = v.Tags
		return nil
	}

	return json.Unmarshal(b, &t.Name)
}
