package router

import (
	"devcode.xeemore.com/systech/gojunkyard/router/internal/httprouter"
	"devcode.xeemore.com/systech/gojunkyard/router/internal/muxie"
)

// EngineType ...
type EngineType uint8

const (
	// HTTPRouter ...
	HTTPRouter EngineType = iota + 1
	// Muxie ...
	Muxie
)

func getRouterEngine(et EngineType) engine {
	switch et {
	case HTTPRouter:
		return httprouter.New()
	case Muxie:
		return muxie.New()
	default:
		panic("Router engine is not found")
	}
}
