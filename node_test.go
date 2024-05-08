package pubsub

import (
	"fmt"
	"sync"
	"testing"
)

func TestNode(t *testing.T) {
	broker := NewBroker()

	first := NewNode(broker)
	second := NewNode(broker)
	third := NewNode(broker)

	if first.HasNext() {
		t.Errorf("Node 1 has no next sibling node")
	}

	first.LinkTo(second)

	if !first.HasNext() {
		t.Errorf("Node 1 has next sibling node")
	}

	first.LinkTo(third)

	if first.next == second {
		t.Errorf("Node 1 is node longer linked to node 2")
	}

	if first.next != third {
		t.Errorf("Node 1 is linked to node 3")
	}

	if third.First() != first {
		t.Errorf("Node 3 should stem from node 1")
	}

	third.LinkTo(second)

	if second.First() != first {
		t.Errorf("Node 2 should stem from node 1")
	}

	if first.Last() != second {
		t.Errorf("Last node for node 1 should be node 2")
	}

	if first.First() != first {
		t.Errorf("If node 1 is the first it should return itself")
	}

	if second.Last() != second {
		t.Errorf("If node 2 is the last it should return itself")
	}

	first.First().UnlinkAll()

	if first.First().HasNext() {
		t.Errorf("First node has no sibling nodes currently")
	}
}

func TestLocalEvent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	broker := NewBroker()

	localEvent := Event(1)

	emitter := NewNode(broker)
	first := NewNode(broker)
	second := NewNode(broker)
	third := NewNode(broker)

	emitter.LinkTo(first)

	emitter.Passthrough(localEvent, emitter.next)

	first.On(localEvent, func() {
		first.LinkTo(second)

		first.Passthrough(localEvent, first.next)

		wg.Done()
	})

	second.On(localEvent, func() {
		second.LinkTo(third)

		wg.Done()
	})
	wg.Wait()

	if first.next != second {
		t.Errorf("Local event not received in first node")
	}

	if second.next != third {
		t.Errorf("Local event not received in second node")
	}
}

func TestGlobalEvent(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	broker := NewBroker()

	event := Event(1)

	emitter := NewNode(broker)

	one := NewNode(broker)
	two := NewNode(broker)
	three := NewNode(broker)
	four := NewNode(broker)
	five := NewNode(broker)

	emitter.LinkTo(one).LinkTo(two).LinkTo(three).LinkTo(four).LinkTo(five)

	one.On(event, func() {
		fmt.Println("received event locally in node one")
	})

	three.OnGlobal(event, func() {
		fmt.Println("received global event in node three")
		wg.Done()
	})

	five.OnGlobal(event, func() {
		fmt.Println("received global event in node five")
		wg.Done()
	})
	emitter.Passthrough(event, emitter.next)
	wg.Wait()
}

func TestUnlinkNext(t *testing.T) {
	broker := NewBroker()

	one, two, three := NewNode(broker), NewNode(broker), NewNode(broker)

	one.LinkTo(two).LinkTo(three)

	one.UnlinkNext()

	if one.next != three {
		t.Errorf("ulink next fails")
	}
}

func TestUnlinkPrev(t *testing.T) {
	broker := NewBroker()

	one, two, three := NewNode(broker), NewNode(broker), NewNode(broker)

	one.LinkTo(two).LinkTo(three)

	three.UnlinkPrev()

	if three.prev != one {
		t.Errorf("ulink prev fails")
	}
}

func TestUnlinkAllPrev(t *testing.T) {
	broker := NewBroker()

	one, two, three, four, five := NewNode(broker), NewNode(broker), NewNode(broker), NewNode(broker), NewNode(broker)

	one.LinkTo(two).LinkTo(three).LinkTo(four).LinkTo(five)

	five.UnlinkAllPrev()

	if five.prev != nil {
		t.Errorf("ulink prev all fails")
	}

	if five.First() != five {
		t.Errorf("node five should be firts")
	}
}

func TestUnlinkAll(t *testing.T) {
	broker := NewBroker()

	one, two, three, four, five := NewNode(broker), NewNode(broker), NewNode(broker), NewNode(broker), NewNode(broker)

	one.LinkTo(two).LinkTo(three).LinkTo(four).LinkTo(five)

	three.UnlinkAll()

	if three.Last() != three {
		t.Errorf("node three should be last")
	}

	if three.First() != three {
		t.Errorf("node three should be first")
	}
}
