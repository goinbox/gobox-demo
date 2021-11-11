package cache

import "gdemo/pcontext"

type Factory interface {
	CacheLogic(ctx *pcontext.Context) Logic
}

var DefaultFactory Factory = new(factory)

type factory struct {
}

func (f *factory) CacheLogic(ctx *pcontext.Context) Logic {
	return NewLogic(ctx)
}
