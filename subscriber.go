package pubsub

type Subscriber struct {
	*Broker
	queue map[string][]Message
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
}

func NewSubscriber(b *Broker) *Subscriber {
	return &Subscriber{b, map[string][]Message{}}
}

func (s *Subscriber) Subscribe(topic string, fn func(m Message)) {
	s.addSubscriber(topic, s, fn)
}

// func (s *Subscriber) Publish(topic string, m Message) {
// 	s.broker.Publish(topic, m)
// }

func (s *Subscriber) Unsubscribe(topic string) {
	s.removeSubscriber(topic, s)
}
