package pipeliner

import "sync"

type pipelinerCmd struct {
	v     interface{}
	resCh chan error
}

var pipelinerCmdPool sync.Pool

func getPipelinerCmd() *pipelinerCmd {
	cmd, _ := pipelinerCmdPool.Get().(*pipelinerCmd)
	if cmd != nil {
		return cmd
	}
	return &pipelinerCmd{resCh: make(chan error, 1)}
}

func putPipelinerCmd(cmd *pipelinerCmd) {
	pipelinerCmdPool.Put(cmd)
}
