package main

import (
	"github.com/bytearena/ecs"
)

type stats struct {
	MaxHealth   int32
	Damage      int32
	Cooldown    float32 `imgui:"%.1f ms"`
	StaminaCost int32
	Dodge       int32
	Heal        int32
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
		MaxHealth:   500,
		Damage:      90,
		Cooldown:    4000,
		StaminaCost: 90,
		Dodge:       10,
	})
	return e
}
