package pubsub

func NewNode(b *Broker) *Node {
	n := &Node{nil, nil, NewListener(), b, NewSubscriber(b), make(chan struct{}, 1)}
	n.watch()
	return n
}

type Node struct {
	next, prev *Node

	*Listener

	*Broker
	*Subscriber
	closeChan chan struct{}
}

func (n *Node) Last() *Node {
	if n.next != nil {
		return n.next.Last()
	}

	return n
}

func (n *Node) HasNext() bool {
	return n.next != nil
}

func (n *Node) First() *Node {
	target := n

	for target.prev != nil {
		target = target.prev
	}

	return target
}

func (n *Node) UnlinkAllPrev() {
	node := n.First()

	if node == n {
		return
	}

	for node != n {

		next := node.next
		node.closeChan <- struct{}{}
		node = next
	}

	n.prev = nil
}

func (n *Node) UnlinkAllNext() {
	last := n.Last()

	for last != n {
		last = last.prev
		last.next.closeChan <- struct{}{}
		last.next = nil

	}

	n.next = nil
}

func (n *Node) UnlinkAll() {
	n.UnlinkAllPrev()
	n.UnlinkAllNext()
}

func (n *Node) UnlinkNext() {
	if n.HasNext() {

		n.next.closeChan <- struct{}{}

		if n.next.HasNext() {

			newNext := n.next.next
			n.next = nil
			n.next = newNext

		}

	}
}

func (n *Node) UnlinkPrev() {
	if n.prev != nil {

		n.prev.closeChan <- struct{}{}

		if n.prev.prev != nil {

			newPrev := n.prev.prev
			n.prev = newPrev
			n.prev.next = n

		}

	}
}

func (n *Node) Next() *Node {
	return n.next
}

func (n *Node) Prev() *Node {
	return n.prev
}

func (n *Node) Unlink() {
	if n.prev != nil {
		n.prev.next.closeChan <- struct{}{}
		n.prev.next = nil
	}
}

func (n *Node) LinkTo(next *Node) *Node {
	n.next = next
	next.prev = n

	return next
}

func (n *Node) watch() {
	go func() {
		for {
			select {

			case <-n.closeChan:
				close(n.closeChan)

				return

			case e := <-n.eventChan:

				switch true {

				case n.HasGlobal(e):

					n.globalEvents[e]()
					n.Passthrough(e, n.next)

				case n.HasLocal(e):

					if n.HasNext() {
						n.Passthrough(e, n.next)
					} else {
						n.events[e]()
					}

				default:

					n.Passthrough(e, n.next)

				}
			}
		}
	}()
}
