package ipmorph

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertNetIPToDecimal(t *testing.T) {
	ipDecimal, err := ConvertNetIPToDecimal("192.168.0.8")
	if assert.NoError(t, err) {
		assert.NotEmpty(t, ipDecimal)
		assert.Equal(t, ipDecimal, uint32(3232235528))
	}
	ipDecimalError, err := ConvertNetIPToDecimal("192.168.0.eight")
	if assert.Error(t, err) {
		assert.Equal(t, ipDecimalError, uint32(0))
		assert.Equal(t, err, errors.New("Invalid ip address value"))
	}
}

func TestConvertDecimalToNetIP(t *testing.T) {
	ipNetIP := ConvertDecimalToNetIP(uint32(3232235528))
	assert.NotEmpty(t, ipNetIP)
	assert.Equal(t, ipNetIP, net.IPv4(192, 168, 0, 8))
}
