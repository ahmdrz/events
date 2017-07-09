package events

import (
	"errors"
	"sync"
	"time"
)

// Handler is main handler for subscribers
type handler func(v interface{}, t time.Time)

type handlers []handler

type _waitList struct {
	wait     sync.WaitGroup
	handlers handlers
}

type _topic struct {
	sync.Mutex
	list map[string]_waitList
}

var topic _topic

func execHandler(wait *sync.WaitGroup, h handler, value interface{}) {
	h(value, time.Now())
	wait.Done()
}

func init() {
	topic.list = make(map[string]_waitList)
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
	topic.Lock()
	defer topic.Unlock()
	h, ok := topic.list[name]
	if !ok {
		return false
	}
	h.wait = sync.WaitGroup{}
	for _, handler := range h.handlers {
		h.wait.Add(1)
		go execHandler(&h.wait, handler, value)
	}
	h.wait.Wait()
	return true
}

// Wait : wait for running handlers
func Wait(name string) {
	h, ok := topic.list[name]
	if ok {
		h.wait.Wait()
	}
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
	topic.list[name] = _waitList{
		handlers: _handlers,
	}
}
