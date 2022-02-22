package pipeliner

import (
	"testing"
)

func Test_getPutPipelinerCmd(t *testing.T) {
	cmd01 := getPipelinerCmd()
	putPipelinerCmd(cmd01)
}
