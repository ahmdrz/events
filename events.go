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
	name  string
}

var topic struct {
	sync.Mutex
	list map[string]waitList
}

func NewRoutineController(name string) *RoutineController {
	return &RoutineController{name: name, group: sync.WaitGroup{}}
}

func (r *RoutineController) PublishRoutine(name string, value interface{}) bool {
	defer r.group.Done()
	if r.name != name {
		return false
	}
	ext := Publish(name, value)
	return ext
}

func (r *RoutineController) AddRoutine(delta int) {
	r.group.Add(delta)
}

func (r *RoutineController) WaitFinish() {
	r.group.Wait()
}

func execHandler(wl *waitList, wg *sync.WaitGroup, h handler, value interface{}) {
	atomic.AddInt32(&wl.runningHandlers, 1)
	go func() {
		defer func() {
			recover()
			atomic.AddInt32(&wl.runningHandlers, -1)
			wg.Done()
		}()
		h(value, time.Now())
	}()
}

func init() {
	topic.list = make(map[string]waitList)
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
	atomic.AddInt32(&h.runningHandlers, 1)
	wait := sync.WaitGroup{}
	for _, handler := range h.handlers {
		wait.Add(1)
		execHandler(&h, &wait, handler, value)
	}
	wait.Wait()
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

// Wait : wait for topic messages
func Wait(name string) {
	h, ok := topic.list[name]
	if !ok {
		return
	}
	for atomic.LoadInt32(&h.runningHandlers) != 0 {
	}
}
