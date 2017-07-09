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

	Wait("user/test")

	if !done {
		t.Fail()
	}
}
