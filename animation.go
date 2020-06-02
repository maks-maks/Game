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
	Speed  float32
}
type AnimationSystem struct{}

func (s *AnimationSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(animationC))
	for _, item := range query {
		animation := item.Components[animationC].(*AnimationComponent)
		animation.Object.Update(time.Now())
	}
}
