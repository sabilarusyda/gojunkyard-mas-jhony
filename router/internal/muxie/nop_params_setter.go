package muxie

type nopParamSetter struct{}

func (*nopParamSetter) Set(string, string) {}
