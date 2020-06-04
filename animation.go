package main

import (
	"time"

	"github.com/bytearena/ecs"
)

type Animatable interface {
	Update(time.Time)
}

type AnimationComponent struct {
	Object Animatable
}

type AnimationSystem struct {
	animations map[string]Animatable
}

func newAnimationSystem() *AnimationSystem {
	return &AnimationSystem{
		animations: make(map[string]Animatable),
	}
}
func (s *AnimationSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(animationC))
	for _, item := range query {
		animation := item.Components[animationC].(*AnimationComponent)
		animation.Object.Update(time.Now())
	}
}
func (s *AnimationSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			damager := ecsManager.GetEntityByID(event.DamagerID, animationC)
			if damager == nil {
				return true
			}
			damagerAnimation := damager.Components[animationC].(*AnimationComponent)

			damagerAnimation.Object = s.animations["tank/attack"]
		}
		return true
	})
}

func (s *AnimationSystem) AddAnimation(key string, a Animatable) {
	s.animations[key] = a
}
