package nsqconsumer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert.Equal(t, &Config{nil, nil, make(consumerConfigs, 0)}, NewConfig(nil, nil))
	assert.Equal(t, &Config{
		[]string{"127.0.0.1:4150"}, []string{"127.0.0.1:4160"}, make(consumerConfigs, 0),
	}, NewConfig([]string{"127.0.0.1:4150"}, []string{"127.0.0.1:4160"}))
}

func TestConfig_AddConsumerConfig(t *testing.T) {
	consumerconfig := NewConsumerConfig("TEST", []Topic{{Name: "TEST_TOPIC"}}, 100, 200, true, true)

	config := NewConfig(nil, nil)
	config.AddConsumerConfig(consumerconfig)

	assert.Equal(t, &Config{nil, nil, consumerConfigs{consumerconfig}}, config)
}

func TestNewConsumerConfig(t *testing.T) {
	var (
		got  = NewConsumerConfig("TEST", []Topic{{Name: "TEST_TOPIC_1"}, {Name: "TEST_TOPIC_2"}}, 100, 200, true, true)
		want = &ConsumerConfig{
			Name:           "TEST",
			Topics:         []Topic{{Name: "TEST_TOPIC_1"}, {Name: "TEST_TOPIC_2"}},
			MaxInFlight:    100,
			Concurrency:    200,
			SkipValidation: true,
			Deduplicator:   true,
		}
	)
	assert.Equal(t, want, got)
}

func Test_consumerConfigs_Decode(t *testing.T) {
	{
		var cc consumerConfigs
		err := cc.Decode(`{}`)
		assert.NotNil(t, err)
	}
	{
		var cc consumerConfigs
		err := cc.Decode(`[{"name": "TEST", "topics": ["TEST_TOPIC_1", "TEST_TOPIC_2"], "maxInFlight": 100, "concurrency": 200, "skipValidation": true, "deduplicator": true}]`)
		assert.Nil(t, err)
		assert.Equal(t, consumerConfigs{NewConsumerConfig("TEST", []Topic{{Name: "TEST_TOPIC_1"}, {Name: "TEST_TOPIC_2"}}, 100, 200, true, true)}, cc)
	}
	{
		var cc consumerConfigs
		err := cc.Decode(`[{"name": "TEST", "topics": [{"name": "TEST_TOPIC_1", "tags": {"source": "registration"}}, "TEST_TOPIC_2"], "maxInFlight": 100, "concurrency": 200, "skipValidation": true, "deduplicator": true}]`)
		assert.Nil(t, err)
		assert.Equal(t, consumerConfigs{NewConsumerConfig("TEST", []Topic{{Name: "TEST_TOPIC_1", Tags: map[string]interface{}{
			"source": "registration",
		}}, {Name: "TEST_TOPIC_2"}}, 100, 200, true, true)}, cc)
	}
}

func TestNewTopic(t *testing.T) {
	type args struct {
		name string
		tags map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *Topic
	}{
		{
			name: "Empty",
			args: args{},
			want: new(Topic),
		},
		{
			name: "Filled",
			args: args{
				name: "ACCOUNTS_REGISTER_DEVICE_IDENTITY",
				tags: map[string]interface{}{
					"TOPIC_NAME": "ACCOUNTS_REGISTER_DEVICE_IDENTITY",
				},
			},
			want: &Topic{
				Name: "ACCOUNTS_REGISTER_DEVICE_IDENTITY",
				Tags: map[string]interface{}{
					"TOPIC_NAME": "ACCOUNTS_REGISTER_DEVICE_IDENTITY",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTopic(tt.args.name, tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
