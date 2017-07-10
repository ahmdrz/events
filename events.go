package events

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Handler is main handler for subscribers
type handler func(v interface{}, t time.Time)

type handlers []handler

type waitList struct {
	wait            sync.WaitGroup
	handlers        handlers
	runningHandlers int32
}

type RoutineController struct {
	group sync.WaitGroup
}

func (r *RoutineController) AddRoutine(delta int) {
	r.group.Add(delta)
}

var routines RoutineController

var topic struct {
	sync.Mutex
	list map[string]waitList
}

func execHandler(wl *waitList, h handler, value interface{}) {
	atomic.AddInt32(&wl.runningHandlers, 1)
	go func() {
		topic.Lock()
		defer topic.Unlock()
		defer func() {
			recover()
			atomic.AddInt32(&wl.runningHandlers, -1)
			wl.wait.Done()
		}()
		h(value, time.Now())
	}()
}

func init() {
	topic.list = make(map[string]waitList)
	routines.group = sync.WaitGroup{}
}

// Close : close('user:login')
func Close(name string) error {
	topic.Lock()
	defer topic.Unlock()
	if _, ok := topic.list[name]; !ok {
		return errors.New("Topic has been closed")
	}
	delete(topic.list, name)
	return nil
}

// Publish : publish('user:login',0)
func Publish(name string, value interface{}) bool {
	defer routines.group.Done()
	h, ok := topic.list[name]
	if !ok {
		return false
	}
	atomic.AddInt32(&h.runningHandlers, 1)
	h.wait = sync.WaitGroup{}
	for _, handler := range h.handlers {
		h.wait.Add(1)
		execHandler(&h, handler, value)
	}
	h.wait.Wait()
	return true
}

// Subscribe : subscribe('user:login')
func Subscribe(name string, h handler) {
	topic.Lock()
	defer topic.Unlock()
	_handlers := make(handlers, 0)
	_, ok := topic.list[name]
	if ok {
		_handlers = topic.list[name].handlers
	}
	_handlers = append(_handlers, h)
	topic.list[name] = waitList{
		handlers: _handlers,
	}
}

func Wait(name string) {
	h, ok := topic.list[name]
	if !ok {
		return
	}
	routines.group.Wait()
	for atomic.LoadInt32(&h.runningHandlers) != 0 {
	}
}
