package events

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

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

func TestPublishers(t *testing.T) {
	testCasesLength := 100 + rand.Intn(200)
	var numbers int
	Subscribe("user:number", func(v interface{}, _t time.Time) {
		numbers++
	})

	for i := 0; i < testCasesLength; i++ {
		Publish("user:number", i)
	}

	Wait("user:number")

	Close("user:number")

	t.Log(numbers, testCasesLength)
	if numbers != testCasesLength {
		t.Fail()
	}
}

func TestSubscribers(t *testing.T) {
	testCasesLength := 100 + rand.Intn(200)
	var numbers int

	for i := 0; i < testCasesLength; i++ {
		Subscribe("user:number", func(v interface{}, _t time.Time) {
			numbers++
		})
	}

	if len(topic.list["user:number"].handlers) != testCasesLength {
		t.Log("handlers problem")
		t.Fail()
	}

	Publish("user:number", 1)

	Wait("user:number")

	Close("user:number")

	t.Log(numbers, testCasesLength)
	if numbers != testCasesLength {
		t.Fail()
	}
}

func TestConcurrentPublishers(t *testing.T) {
	testCasesLength := 100 + rand.Intn(200)
	var numbers int
	Subscribe("user:number", func(v interface{}, _t time.Time) {
		numbers++
	})

	for i := 0; i < testCasesLength; i++ {
		go Publish("user:number", i)
	}

	Wait("user:number")

	Close("user:number")

	t.Log(numbers, testCasesLength)
	if numbers != testCasesLength {
		t.Fail()
	}
}
