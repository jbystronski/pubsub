package pubsub

type Subscriber struct {
	broker *Broker
	queue  map[string][]Message
}

func (s *Subscriber) EnqueueMessage(label string, m Message) {
	if _, ok := s.queue[label]; !ok {
		s.queue[label] = []Message{}
	}

	s.queue[label] = append(s.queue[label], m)
}

func (s *Subscriber) DequeueMessage(label string) Message {
	var m Message

	if _, ok := s.queue[label]; ok {
		if len(s.queue[label]) > 0 {

			m = s.queue[label][0]
			s.queue[label] = s.queue[label][1:]
		}
	}
	return m
	// m = s.Queue[label][0]

	// if len(s.Queue[label]) > 0 {
	// 	s.Queue[label] = s.Queue[label][1:]
	// } else {
	// 	delete(s.Queue, label)
	// }

	// return m
}

func NewSubscriber(b *Broker) *Subscriber {
	return &Subscriber{broker: b, queue: map[string][]Message{}}
}

func (s *Subscriber) Subscribe(topic string, fn func(m Message)) {
	s.broker.addSubscriber(topic, s, fn)
}

func (s *Subscriber) Unsubscribe(topic string) {
	s.broker.removeSubscriber(topic, s)
}
