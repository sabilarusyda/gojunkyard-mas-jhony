package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	assert.Len(t, uuid, 32)
	assert.NotEmpty(t, uuid)
	assert.NotContains(t, uuid, "-")
}

func Benchmark_GenerateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUUID()
	}
}
