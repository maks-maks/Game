package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventBus(t *testing.T) {
	b := NewEventBus()
	b.Schedule(&TestEvent{V: 1}, 0)
	b.Schedule(&TestEvent{V: 2}, 100)
	b.Schedule(&TestEvent{V: 3}, 0)
	b.Schedule(&TestEvent{V: 4}, 50)

	actual := []float32{
		heap.Pop(b.scheduled).(*TimedEvent).Time,
		heap.Pop(b.scheduled).(*TimedEvent).Time,
		heap.Pop(b.scheduled).(*TimedEvent).Time,
		heap.Pop(b.scheduled).(*TimedEvent).Time,
	}
	assert.Equal(t, []float32{0, 0, 50, 100}, actual)
}

type TestEvent struct {
	V int
}
