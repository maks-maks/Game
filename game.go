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
	ecsManager.RegisterComponent("target", &TargetComponent{})

	createTank("Frederik")
	createTank("Frederik2")

	// e2 := ecsManager.NewEntity()
	// ecsManager.AddComponent(e2, &PositionComponent{X: 10, Y: 15})
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
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}

type TargetComponent struct {
	TargetID ecs.EntityID
}

type targetingSystem struct{}

func (s *targetingSystem) Update(dt float32) {
	entities := ecsManager.Query(ecs.BuildTag(ecsManager.componentMap["target"])).Entities()

	for _, e := range entities {
		d, _ := e.GetComponentData(ecsManager.componentMap["target"])
		data := d.(*TargetComponent)

		targets := ecsManager.Query(ecs.BuildTag(ecsManager.componentMap["stats"])).Entities()
		for _, t := range targets {
			if t.ID != e.ID {
				data.TargetID = t.ID
			}
		}
	}
}
