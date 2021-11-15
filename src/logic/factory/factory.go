package factory

import (
	"gdemo/logic/cache"
)

type Factory interface {
	CacheLogic() cache.Logic
}

var DefaultFactory Factory = new(factory)

type factory struct {
}

func (f *factory) CacheLogic() cache.Logic {
	return cache.NewLogic()
}
