package main

import (
	"time"

	"github.com/bytearena/ecs"
)

type stats struct {
	maxHealth   int
	damage      int
	cooldown    time.Duration
	staminaCost int
	dodge       int
	heal        int
}

func setupECS() {
	ecsManager = NewECSManager()

	ecsManager.RegisterComponent("position", &PositionComponent{})
	ecsManager.RegisterComponent("stats", &stats{})

	createTank("Frederik")

	e2 := ecsManager.NewEntity()
	ecsManager.AddComponent(e2, &PositionComponent{X: 10, Y: 15})
}

func createTank(n string) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &PositionComponent{X: 1, Y: 2})
	ecsManager.AddComponent(e, &stats{
		maxHealth:   500,
		damage:      90,
		cooldown:    4000 * time.Millisecond,
		staminaCost: 90,
		dodge:       10,
	})
	return e
}
