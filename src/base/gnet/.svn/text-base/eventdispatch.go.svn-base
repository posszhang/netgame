package gnet

import ()

type EventFunc func(data interface{})

type EventDispatch struct {
	eventMap map[string]EventFunc
}

func (ev *EventDispatch) AddEventListener(name string, fun EventFunc) {

	if ev.eventMap == nil {
		ev.eventMap = make(map[string]EventFunc)
	}

	ev.eventMap[name] = fun
}

func (ev *EventDispatch) RemoveEventListener(name string, fun EventFunc) {
	if _, ok := ev.eventMap[name]; !ok {
		return
	}

	delete(ev.eventMap, name)
}

func (ev *EventDispatch) DispatchEvent(name string, data interface{}) {

	fun, ok := ev.eventMap[name]
	if !ok {
		return
	}

	fun(data)
}
