package pubsub

import (
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	message := Message("Checking broker functionallity")

	broker := NewBroker()

	s1 := NewSubscriber(broker)
	s2 := NewSubscriber(broker)

	broker.Publish("test_topic", message)

	s1.Subscribe("test_topic", func(m Message) {
		want := string(message)

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s2.Subscribe("test_topic", func(m Message) {
		want := string(message)

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s2.Subscribe("message from first", func(m Message) {
		want := string("Checking communiaction with subscriber 1")

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s1.Subscribe("message from second", func(m Message) {
		want := string("Checking communiaction with subscriber 2")

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s1.Publish("message from first", Message("Checking communiaction with subscriberssssss 1"))
	s2.Publish("message from", Message("Checking communiaction with subscriber 2"))

	time.Sleep(time.Second * 2)
}
