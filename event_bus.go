package main

import (
	"fmt"
	"time"
)

type EventBus struct {
	queue     []Event
	scheduled []TimedEvent
}

func (b *EventBus) Publish(e Event) {
	log = append(log, fmt.Sprintf("Event %s", e))
	b.queue = append(b.queue, e)
}

func (b *EventBus) Schedule(e Event, dt int) {
	// find an index to insert
	insertIndex := 0
	t := time.Now().Add(time.Duration(dt) * time.Millisecond)
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

func (b *EventBus) AdvanceScheduled(t time.Time) {
	for _, item := range b.scheduled {
		if item.Time.Before(t) {
			log = append(log, fmt.Sprintf("%s", b.scheduled[0].Event))
			b.queue = append(b.queue, b.scheduled[0].Event)
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
	Event Event
	Time  time.Time
}
