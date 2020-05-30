package main

import (
	"time"
)

type EventBus struct {
	queue     []Event
	scheduled []TimedEvent
}

func (b *EventBus) Publish(e Event) {
	b.queue = append(b.queue, e)
}

func (b *EventBus) Schedule(e Event, t time.Time) {
	// find an index to insert
	insertIndex := 0
	for i, item := range b.scheduled {
		if item.Time.After(t) {
			insertIndex = i
			break
		}
	}

	// insert at proper place
	b.scheduled = append(b.scheduled, TimedEvent{})
	copy(b.scheduled[insertIndex+1:], b.scheduled[insertIndex:])
	b.scheduled[insertIndex] = TimedEvent{
		Time:  t,
		Event: e,
	}
}

func (b *EventBus) AdvanceScheduled() {
	for _, item := range b.scheduled {
		if item.Time.Before(time.Now()) {
			b.queue = append(b.queue, b.scheduled[0])
			b.scheduled = b.scheduled[1:]
		} else {
			break
		}
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
	Time  time.Time
	Event Event
}
