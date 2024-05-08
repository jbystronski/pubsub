package subscriber_test

import (
	"testing"
	"time"

	"github.com/jbystronski/pubsub"
)

func TestSubscribe(t *testing.T) {
	message := pubsub.Message("Checking broker functionallity")

	broker := pubsub.NewBroker()

	s1 := pubsub.NewSubscriber(broker)
	s2 := pubsub.NewSubscriber(broker)

	broker.Publish("test_topic", message)

	s1.Subscribe("test_topic", func(m pubsub.Message) {
		want := string(message)

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s2.Subscribe("test_topic", func(m pubsub.Message) {
		want := string(message)

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s2.Subscribe("message from first", func(m pubsub.Message) {
		want := string("Checking communiaction with subscriber 1")

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s1.Subscribe("message from second", func(m pubsub.Message) {
		want := string("Checking communiaction with subscriber 2")

		got := string(m)

		if want != got {
			t.Errorf("Received: %s; wanted %s", got, want)
		}
	})

	s1.Publish("message from first", pubsub.Message("Checking communiaction with subscriberssssss 1"))
	s2.Publish("message from", pubsub.Message("Checking communiaction with subscriber 2"))

	time.Sleep(time.Second * 2)
}
