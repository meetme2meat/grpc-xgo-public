package model

import "xgo/main/src/event"

type HookFunc func(event.EventRecord)

var hooks = make(map[string]HookFunc)

func AddHook(name string, fn HookFunc) {
	hooks[name] = fn
}

func RemoveHook(name string) {
	delete(hooks, name)
}
