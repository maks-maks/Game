package main

import (
	"container/heap"
	"fmt"
	"time"
)

type EventBus struct {
	queue     []Event
	scheduled *EventPriorityQueue
	curTime   float32
}

func NewEventBus() *EventBus {
	epq := &EventPriorityQueue{}
	heap.Init(epq)

	return &EventBus{
		scheduled: epq,
	}
}

func (b *EventBus) Publish(e Event) {
	b.Schedule(e, 0)
}

func (b *EventBus) Schedule(e Event, dt float32) {
	heap.Push(b.scheduled, &TimedEvent{
		Time:  b.curTime + dt,
		Event: e,
	})
	log = append(log, fmt.Sprintf("Scheduled after %v %s", time.Duration(int(dt))*time.Millisecond, e))
}

func (b *EventBus) AdvanceScheduled(dt float32) {
	b.curTime += dt
	for len(*b.scheduled) > 0 {
		if (*b.scheduled)[0].Time > b.curTime {
			break
		}

		v := heap.Pop(b.scheduled).(*TimedEvent)
		log = append(log, fmt.Sprintf("%s", v.Event))
		b.queue = append(b.queue, v.Event)
	}
}

func (b *EventBus) Iterate(f func(Event) bool) {
	for _, e := range b.queue {
		f(e)
		// TODO remove from queue imidiately if f returns false
	}
}

func (b *EventBus) ClearQueue() {
	b.queue = b.queue[:0]
}

type Event interface{}

type TimedEvent struct {
	Event Event
	Time  float32
}

type EventPriorityQueue []*TimedEvent

func (pq EventPriorityQueue) Len() int {
	return len(pq)
}

func (pq EventPriorityQueue) Less(i, j int) bool {
	return pq[i].Time <= pq[j].Time
}

func (pq EventPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *EventPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*TimedEvent))
}

func (pq *EventPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
