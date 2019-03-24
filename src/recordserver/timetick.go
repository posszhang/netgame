package main

import (
	"base/util"
)

type TimeTick struct {
	util.TimeTick
}

func NewTimeTick() *TimeTick {
	tick := &TimeTick{}
	tick.Derived = tick
	tick.Run()

	return tick
}

func (tt *TimeTick) TimeAction() {
	routeManager.Do()
}
