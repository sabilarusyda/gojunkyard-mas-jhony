package nop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNopReporter(t *testing.T) {
	assert.Equal(t, new(Nop), NewNopReporter())
}

func TestNopReporter(t *testing.T) {
	nop := NewNopReporter()
	nop.Debug("")
	nop.Debugf("", "")
	nop.Debugln("")
	nop.Info("")
	nop.Infof("", "")
	nop.Infoln("")
	nop.Warning("")
	nop.Warningf("", "")
	nop.Warningln("")
	nop.Error("")
	nop.Errorf("", "")
	nop.Errorln("")
	assert.Nil(t, nop.ReportPanic(nil, nil))
	assert.Nil(t, nop.ReportHTTPPanic(nil, nil, nil))
}
