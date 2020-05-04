package main

import (
	"math/rand"

	"github.com/bytearena/ecs"
)

type StatsComponent struct {
	MaxHealth   int32
	Damage      int32
	Cooldown    float32 `imgui:"%.1f ms"`
	StaminaCost int32
	Stamina     int32
	Dodge       int32
	Heal        int32
	Health      int32
	Reload      float32 `imgui:"%.1f ms"`
}

func setupECS() {
	ecsManager = NewECSManager()

	ecsManager.RegisterComponent("position", &PositionComponent{})
	ecsManager.RegisterComponent("stats", &StatsComponent{})
	ecsManager.RegisterComponent("target", &TargetComponent{})
	for i := 1; i < 4; i++ {
		createTank("Frederik")
	}

	// e2 := ecsManager.NewEntity()
	// ecsManager.AddComponent(e2, &PositionComponent{X: 10, Y: 15})
}

func createTank(n string) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &PositionComponent{
		X: (rand.Float32() * 300) + 100,
		Y: (rand.Float32() * 200) + 100,
	})
	ecsManager.AddComponent(e, &StatsComponent{
		MaxHealth:   500,
		Health:      500,
		Damage:      90,
		Cooldown:    1000,
		Stamina:     100,
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

type battleSystem struct{}

func (s *battleSystem) Update(dt float32) {
	targetC := ecsManager.componentMap["target"]
	statC := ecsManager.componentMap["stats"]

	query := ecsManager.Query(ecs.BuildTag(targetC, statC))

	for _, item := range query {
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statC].(*StatsComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, statC)
		if target == nil {
			currentTarget.TargetID = 0
			continue
		}

		targetStats := target.Components[statC].(*StatsComponent)

		if stats.Reload < stats.Cooldown {
			stats.Reload = stats.Reload + dt
			continue
		}
		if stats.Stamina < stats.StaminaCost {
			stats.Stamina = stats.Stamina + 50
			stats.Reload = 0
			continue
		}

		targetStats.Health = targetStats.Health - stats.Damage
		stats.Stamina = stats.Stamina - stats.StaminaCost
		stats.Reload = 0
		if targetStats.Health <= 0 {
			ecsManager.DisposeEntity(target.Entity)
		}
	}
}
