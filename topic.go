package pubsub

type topic struct {
	done        chan struct{}
	messageChan chan Message

	subscribers map[*Subscriber]func(Message)
}

// TODO: make subscribers into a sync map
func newTopic() *topic {
	t := topic{
		messageChan: make(chan Message, 1),
		subscribers: map[*Subscriber]func(Message){},
	}

	return &t
}

func (t *topic) send(msg Message) {
	t.messageChan <- msg
}
