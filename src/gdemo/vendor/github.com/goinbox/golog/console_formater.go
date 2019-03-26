package golog

import (
	"github.com/goinbox/color"
)

type colorFunc func(msg []byte) []byte

type consoleFormater struct {
	f IFormater

	levelColorFuncs map[int]colorFunc
}

func NewConsoleFormater(f IFormater) *consoleFormater {
	c := &consoleFormater{
		f: f,

		levelColorFuncs: map[int]colorFunc{
			LEVEL_DEBUG:     color.Yellow,
			LEVEL_INFO:      color.Blue,
			LEVEL_NOTICE:    color.Cyan,
			LEVEL_WARNING:   color.Maganta,
			LEVEL_ERROR:     color.Red,
			LEVEL_CRITICAL:  color.Black,
			LEVEL_ALERT:     color.White,
			LEVEL_EMERGENCY: color.Green,
		},
	}

	return c
}

func (c *consoleFormater) SetColor(level int, cf colorFunc) *consoleFormater {
	c.levelColorFuncs[level] = cf

	return c
}

func (c *consoleFormater) Format(level int, msg []byte) []byte {
	return c.levelColorFuncs[level](c.f.Format(level, msg))
}
