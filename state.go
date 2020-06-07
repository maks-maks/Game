package main

import (
	"github.com/bytearena/ecs"
)

type StateComponent struct {
	// idle, attack, walk, dodge, dying, dead
	State string
}

type StateSystem struct{}

func (s *StateSystem) Update(dt float32) {}

func (s *StateSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			s.setState(event.DamagerID, "attack")
		case *WalkStartEvent:
			s.setState(event.EntityID, "walk")
		case *StopEvent:
			s.setState(event.EntityID, "idle")
		}
		return true
	})
}

func (s *StateSystem) setState(id ecs.EntityID, state string) bool {
	entity := ecsManager.GetEntityByID(id, stateC)
	if entity == nil {
		return false
	}
	entity.Components[stateC].(*StateComponent).State = state
	return true
}
