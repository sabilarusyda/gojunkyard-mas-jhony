package radix

import (
	"errors"
	"testing"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (cm *ClientMock) Do(action radix.Action) error {
	return cm.Called(action).Error(0)
}

func (cm *ClientMock) Close() error {
	return cm.Called().Error(0)
}

func TestNew(t *testing.T) {
	var client = &radix.Pool{}
	assert.Equal(t, &Radix{client: client}, New(client))
}

func TestRadix_SetNX(t *testing.T) {
	const key = "ahuehue:ahuehue"
	var b = false

	pipe := radix.Pipeline(
		radix.Cmd(&b, "SETNX", key, ""),
		radix.Cmd(nil, "EXPIRE", key, "1"),
	)

	cm := new(ClientMock)
	cm.On("Do", pipe).Return(nil)

	client := New(cm)
	got, err := client.SetNX(key, time.Second)
	assert.Nil(t, err)
	assert.False(t, got)
}

func TestRadix_Expire(t *testing.T) {
	{
		const key = "ahuehue:ahuehue"
		cmd := radix.Cmd(nil, "DEL", key)
		cm := new(ClientMock)
		cm.On("Do", cmd).Return(nil)
		client := New(cm)
		err := client.Delete(key)
		assert.Nil(t, err)
	}
	{
		const key = "ahuehue:ahuehue"
		cmd := radix.Cmd(nil, "DEL", key)
		cm := new(ClientMock)
		cm.On("Do", cmd).Return(errors.New("HEUHEUHEUHE"))
		client := New(cm)
		err := client.Delete(key)
		assert.Equal(t, errors.New("HEUHEUHEUHE"), err)
	}
}
