package pubsub

import "sync"

type Message string

type Broker struct {
	list map[string]*topic
	sync.Mutex
}

var singleBroker *Broker

var once sync.Once

func GlobalBroker() *Broker {
	once.Do(func() {
		singleBroker = NewBroker()
	})

	return singleBroker
}

func NewBroker() *Broker {
	return &Broker{
		list: map[string]*topic{},
	}
}

func (b *Broker) topic(name string) *topic {
	return b.list[name]
}

func (b *Broker) newTopic(name string) {
	b.list[name] = newTopic()

	go func() {
		for {
			select {
			case <-b.topic(name).done:
				return

			case msg := <-b.topic(name).messageChan:

				b.Lock()
				for _, callbackFn := range b.topic(name).subscribers {
					callbackFn(msg)
				}
				b.Unlock()

			}
		}
	}()
}

func (b *Broker) hasTopic(name string) bool {
	if _, ok := b.list[name]; ok {
		return true
	}

	return false
}

func (b *Broker) Close(topic string) {
	if b.hasTopic(topic) {

		b.topic(topic).done <- struct{}{}

		delete(b.list, topic)
	}
}

func (b *Broker) addSubscriber(name string, s *Subscriber, fn func(m Message)) {
	if !b.hasTopic(name) {
		b.newTopic(name)
	}

	b.topic(name).subscribers[s] = fn
}

func (b *Broker) removeSubscriber(name string, s *Subscriber) {
	delete(b.topic(name).subscribers, s)
}

func (b *Broker) Publish(topic string, m Message) {
	if !b.hasTopic(topic) {
		b.newTopic(topic)
	}

	t := b.topic(topic)

	if t != nil {
		t.send(m)
	}
}
