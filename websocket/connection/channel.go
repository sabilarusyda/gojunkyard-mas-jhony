package connection

type channel struct {
	conns connectionSet
}

func newChannel() *channel {
	return &channel{}
}

type channelSet struct {
	m map[*channel]struct{}
}

func (cs *channelSet) add(chans ...*channel) {
	if cs.m == nil {
		cs.m = make(map[*channel]struct{}, len(chans))
	}
	for _, ch := range chans {
		cs.m[ch] = empty
	}
}

func (cs *channelSet) each(f func(c *channel)) {
	for c := range cs.m {
		f(c)
	}
}

func (cs *channelSet) exist(ch *channel) bool {
	_, ok := cs.m[ch]
	return ok
}

func (cs *channelSet) remove(chans ...*channel) {
	for _, ch := range chans {
		delete(cs.m, ch)
	}
}
