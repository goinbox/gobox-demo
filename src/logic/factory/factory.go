package factory

import (
	"gdemo/logic/cache"
	"gdemo/logic/demo"
)

type LogicFactory interface {
	CacheLogic() cache.Logic
	DemoLogic() demo.Logic
}

var DefaultLogicFactory LogicFactory = new(logicFactory)

type logicFactory struct {
}

func (f *logicFactory) CacheLogic() cache.Logic {
	return cache.NewLogic()
}

func (f *logicFactory) DemoLogic() demo.Logic {
	return demo.NewLogic()
}
