package events

import (
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	var published bool
	published = Publish("user/test", 10)
	if published {
		t.Fail()
	}

	Subscribe("user/test", func(v interface{}, t time.Time) {
		// ignore
	})

	published = Publish("user/test", 10)
	if !published {
		t.Fail()
	}

	Close("user/test")
}

func TestSubscribe(t *testing.T) {
	done := false
	Subscribe("user/test", func(v interface{}, t time.Time) {
		done = true
	})

	Publish("user/test", 10)

	Close("user/test")

	if !done {
		t.Fail()
	}
}

func TestMultiEvents(t *testing.T) {
	numbers := 0
	Subscribe("user/login", func(v interface{}, _t time.Time) {
		t.Log("user logged in", v, _t)
		numbers++
	})

	Subscribe("user/logout", func(v interface{}, _t time.Time) {
		t.Log("user logged out", v, _t)
		numbers++
	})

	Publish("user/login", "ahmdrz")
	Publish("user/logout", "ahmdrz")

	Close("user/login")
	Close("user/logout")

	if numbers != 2 {
		t.Fatal("one of subscribers not loaded")
	}
}

func TestConcurrentPublishers(t *testing.T) {
	var numbers int32
	Subscribe("user:number", func(v interface{}, _t time.Time) {
		numbers++
	})
	c := NewRountineControler("user:number")
	c.AddRoutine(1000)
	for i := 0; i < 1000; i++ {
		go c.PublishRoutine("user:number", i)
	}
	c.WaitFinish()

	t.Log(numbers, 1000)
	if numbers != 1000 {
		t.Fail()
	}
}
