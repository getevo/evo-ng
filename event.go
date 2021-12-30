package evo

import (
	"os"
	"os/signal"
	"syscall"
)

type event struct {
	events map[string][]func()
}

var Events = event{}

func (e *event) On(event string, callback func()) {
	if _, ok := e.events[event]; ok {
		e.events[event] = append(e.events[event], callback)
	} else {
		e.events[event] = []func(){callback}
	}
}

func (e *event) Trigger(event string) {
	if callbacks, ok := e.events[event]; ok {
		for _, callback := range callbacks {
			callback()
		}
	}
}

func (e *event) List() []string {
	var keys []string
	for key, _ := range e.events {
		keys = append(keys, key)
	}
	return keys
}

func (e *event) Register() {
	if len(e.events) > 0 {
		return
	}
	e.events = map[string][]func(){}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Events.Trigger("exit")
		os.Exit(1)
	}()
}
