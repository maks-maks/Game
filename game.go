package main

import (
	"math"
	"math/rand"

	"github.com/bytearena/ecs"
)

type Ability interface {
	Activate(id ecs.EntityID)
	Deactivate(id ecs.EntityID)
}

type RageAbility struct{}

func (a *RageAbility) Activate(id ecs.EntityID) {
	statC := ecsManager.componentMap["stats"]

	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 2
	stats.Cooldown = stats.Cooldown / 3
}
func (a *RageAbility) Deactivate(id ecs.EntityID) {
	statC := ecsManager.componentMap["stats"]

	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 1
	stats.Cooldown = stats.Cooldown * 3
}

type UltimateComponent struct {
	Cooldown float32 `imgui:"%.1f ms"`
	Reload   float32 `imgui:"%.1f ms"`
	Ability  Ability
	Active   bool
}

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
	AttackRange float32
	DodgeRange  float32
	Resist      float32
}
type SquadComponent struct {
	Squad string
}

func setupECS() {
	ecsManager = NewECSManager()

	ecsManager.RegisterComponent("position", &PositionComponent{})
	ecsManager.RegisterComponent("stats", &StatsComponent{})
	ecsManager.RegisterComponent("target", &TargetComponent{})
	ecsManager.RegisterComponent("squad", &SquadComponent{})
	ecsManager.RegisterComponent("ulta", &UltimateComponent{})
	//createTank("Frederik", "a", 100, 100)
	//createRanger("Legolas", "b", 400, 400)
	// for i := 1; i < 5; i++ {
	createSquad("Geroi", 400, 400)
	createSquad("Sandali", 100, 100)
	//createSquad("Angels", 400, 100)
	//createSquad("Daemons", 100, 400)
	// createTank("Frederik", "Sandali", 100, 100)
	// createTank("Frederik", "Geroi", 400, 400)
	// createRanger("Legolas", "Sandali", 100, 100)
	// createRanger("Legolas", "Geroi", 400, 400)
	// }

	// e2 := ecsManager.NewEntity()
	// ecsManager.AddComponent(e2, &PositionComponent{X: 10, Y: 15})
}

func createSquad(squad string, x float32, y float32) {
	createHealer("Angel", squad, x, y)
	createRanger("Legolas", squad, x, y)

	createTank("Frederik", squad, x, y)
}
func createHealer(n string, squad string, x float32, y float32) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &PositionComponent{
		X: x + rand.Float32()*200 - 100,
		Y: y + rand.Float32()*200 - 100,
	})
	ecsManager.AddComponent(e, &SquadComponent{
		Squad: squad,
	})
	ecsManager.AddComponent(e, &StatsComponent{
		MaxHealth:   200,
		Health:      200,
		Damage:      20,
		Cooldown:    4000,
		Stamina:     100,
		StaminaCost: 70,
		Dodge:       30,
		AttackRange: 75,
		DodgeRange:  30,
		Heal:        100,
		Resist:      1,
	})
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}
func createRanger(n string, squad string, x float32, y float32) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &PositionComponent{
		X: x + rand.Float32()*200 - 100,
		Y: y + rand.Float32()*200 - 100,
	})
	ecsManager.AddComponent(e, &SquadComponent{
		Squad: squad,
	})
	ecsManager.AddComponent(e, &StatsComponent{
		MaxHealth:   200,
		Health:      200,
		Damage:      20,
		Cooldown:    400,
		Stamina:     100,
		StaminaCost: 20,
		Dodge:       30,
		AttackRange: 300,
		DodgeRange:  25,
		Resist:      1,
	})
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}

func createTank(n string, squad string, x float32, y float32) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &PositionComponent{
		X: x + rand.Float32()*200 - 100,
		Y: y + rand.Float32()*200 - 100,
	})
	ecsManager.AddComponent(e, &SquadComponent{
		Squad: squad,
	})
	ecsManager.AddComponent(e, &UltimateComponent{
		Cooldown: 10000,
		Ability:  &RageAbility{},
	})
	ecsManager.AddComponent(e, &StatsComponent{
		MaxHealth:   500,
		Health:      500,
		Damage:      90,
		Cooldown:    3000,
		Stamina:     100,
		StaminaCost: 90,
		Dodge:       10,
		AttackRange: 75,
		DodgeRange:  -25,
		Resist:      1,
	})
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}

type TargetComponent struct {
	TargetID ecs.EntityID
}

type targetingSystem struct{}

func (s *targetingSystem) Update(dt float32) {
	targetC := ecsManager.componentMap["target"]
	statC := ecsManager.componentMap["stats"]
	positionC := ecsManager.componentMap["position"]
	squadC := ecsManager.componentMap["squad"]

	entities := ecsManager.Query(ecs.BuildTag(targetC, squadC, statC, positionC))

	for _, item := range entities {
		target := item.Components[targetC].(*TargetComponent)
		squad := item.Components[squadC].(*SquadComponent)
		stats := item.Components[statC].(*StatsComponent)
		position := item.Components[positionC].(*PositionComponent)
		targetEntities := ecsManager.Query(ecs.BuildTag(statC, squadC, positionC))

		var curDistance float32 = 1000000

		for _, te := range targetEntities {
			targetSquad := te.Components[squadC].(*SquadComponent)
			targetStats := te.Components[statC].(*StatsComponent)
			targetPosition := te.Components[positionC].(*PositionComponent)

			d := distance(position, targetPosition)
			if d < curDistance {
				if stats.Heal > 0 {
					if te.Entity.ID != item.Entity.ID && squad.Squad == targetSquad.Squad && targetStats.Health < targetStats.MaxHealth {
						curDistance = d
						target.TargetID = te.Entity.ID

					}
				} else if te.Entity.ID != item.Entity.ID && squad.Squad != targetSquad.Squad {
					curDistance = d
					target.TargetID = te.Entity.ID
				}
			}
		}
	}
}

type movementSystem struct{}

func (s *movementSystem) Update(dt float32) {
	targetC := ecsManager.componentMap["target"]
	positionC := ecsManager.componentMap["position"]
	statsC := ecsManager.componentMap["stats"]

	query := ecsManager.Query(ecs.BuildTag(targetC, positionC, statsC))
	for _, item := range query {
		position := item.Components[positionC].(*PositionComponent)
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statsC].(*StatsComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, positionC)
		if target == nil {
			currentTarget.TargetID = 0
			continue
		}

		targetPosition := target.Components[positionC].(*PositionComponent)
		x1, y1, x2, y2 := position.X, position.Y, targetPosition.X, targetPosition.Y

		d := distance(position, targetPosition)

		if d < stats.AttackRange {
			continue
		}

		position.X += (x2 - x1) / d * dt / 30
		position.Y += (y2 - y1) / d * dt / 30
	}

}

func distance(a, b *PositionComponent) float32 {
	d2 := ((b.X - a.X) * (b.X - a.X)) + ((b.Y - a.Y) * (b.Y - a.Y))
	return float32(math.Sqrt(float64(d2)))
}

type ultimatesSystem struct{}

func (s *ultimatesSystem) Update(dt float32) {
	// targetC := ecsManager.componentMap["target"]
	// statC := ecsManager.componentMap["stats"]
	// positionC := ecsManager.componentMap["position"]
	ultaC := ecsManager.componentMap["ulta"]
	query := ecsManager.Query(ecs.BuildTag(ultaC))

	for _, item := range query {
		ulta := item.Components[ultaC].(*UltimateComponent)
		if ulta.Reload < ulta.Cooldown {
			ulta.Reload = ulta.Reload + dt
			continue
		}
		ulta.Reload = 0
		if ulta.Active {
			ulta.Ability.Deactivate(item.Entity.ID)
			ulta.Active = false
		} else {
			ulta.Ability.Activate(item.Entity.ID)
			ulta.Active = true
		}
	}
}

type battleSystem struct{}

func (s *battleSystem) Update(dt float32) {
	targetC := ecsManager.componentMap["target"]
	statC := ecsManager.componentMap["stats"]
	positionC := ecsManager.componentMap["position"]

	query := ecsManager.Query(ecs.BuildTag(targetC, statC, positionC))

	for _, item := range query {
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statC].(*StatsComponent)
		position := item.Components[positionC].(*PositionComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, statC, positionC)
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

		targetPosition := target.Components[positionC].(*PositionComponent)
		x1, y1, x2, y2 := position.X, position.Y, targetPosition.X, targetPosition.Y
		d := distance(position, targetPosition)
		if d > stats.AttackRange {
			continue
		}

		stats.Stamina = stats.Stamina - stats.StaminaCost
		stats.Reload = 0
		if stats.Heal > 0 {
			targetStats.Health = targetStats.Health + stats.Heal
			if targetStats.Health > targetStats.MaxHealth {
				targetStats.Health = targetStats.MaxHealth
				continue
			}
			continue
		}
		if rand.Int31n(100)+1 <= targetStats.Dodge {
			targetPosition.X += (x2 - x1) / d * 1 * targetStats.DodgeRange
			targetPosition.Y += (y2 - y1) / d * 1 * targetStats.DodgeRange
			continue
		}

		targetStats.Health = targetStats.Health - int32((float32(stats.Damage) / targetStats.Resist))
		if targetStats.Health <= 0 {
			ecsManager.DisposeEntity(target.Entity)
		}
	}
}
