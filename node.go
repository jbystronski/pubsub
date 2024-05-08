package pubsub

func NewNode(b *Broker) *Node {
	n := &Node{nil, nil, NewListener(), b, NewSubscriber(b), make(chan struct{}, 1)}
	n.watch()
	return n
}

type Node struct {
	Next, Prev *Node

	*Listener

	*Broker
	*Subscriber
	closeChan chan struct{}
}

func (n *Node) Last() *Node {
	if n.Next != nil {
		return n.Next.Last()
	}

	return n
}

func (n *Node) HasNext() bool {
	return n.Next != nil
}

func (n *Node) First() *Node {
	target := n

	for target.Prev != nil {
		target = target.Prev
	}

	return target
}

func (n *Node) UnlinkAll() {
	last := n.Last()

	for last != n {
		last = last.Prev
		last.Next.closeChan <- struct{}{}
		last.Next = nil

	}

	n.Next = nil
}

func (n *Node) Unlink() {
	if n.Prev != nil {
		n.Prev.Next.closeChan <- struct{}{}
		n.Prev.Next = nil
	}
}

func (n *Node) LinkTo(next *Node) *Node {
	n.Next = next
	next.Prev = n

	return next
}

func (n *Node) watch() {
	go func() {
		for {
			select {

			case <-n.closeChan:
				return

			case e := <-n.eventChan:

				switch true {

				case n.HasGlobal(e):

					n.globalEvents[e]()
					n.Passthrough(e, n.Next)

				case n.HasLocal(e):

					if n.HasNext() {
						n.Passthrough(e, n.Next)
					} else {
						n.events[e]()
					}

				default:

					n.Passthrough(e, n.Next)

				}
			}
		}
	}()
}
