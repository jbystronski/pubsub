package pubsub

type Listener struct {
	events,
	globalEvents map[Event]func()

	eventChan chan Event
}

func NewListener() *Listener {
	return &Listener{map[Event]func(){}, map[Event]func(){}, make(chan Event, 1)}
}

func (l Listener) Passthrough(e Event, n *Node) {
	if n != nil {
		n.eventChan <- e
	}
}

func (l *Listener) HasGlobal(e Event) bool {
	if _, ok := l.globalEvents[e]; ok {
		return true
	}
	return false
}

func (l *Listener) HasLocal(e Event) bool {
	if _, ok := l.events[e]; ok {
		return true
	}
	return false
}

func (l Listener) On(e Event, fn func()) {
	l.events[e] = fn
}

func (l Listener) OnGlobal(e Event, fn func()) {
	l.globalEvents[e] = fn
}
