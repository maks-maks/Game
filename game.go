package main

import (
	"math"
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
	AttackRange float32
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
	createTank("Frederik", squad, x, y)
	createRanger("Legolas", squad, x, y)
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
		Damage:      25,
		Cooldown:    400,
		Stamina:     100,
		StaminaCost: 20,
		Dodge:       30,
		AttackRange: 225,
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
	ecsManager.AddComponent(e, &StatsComponent{
		MaxHealth:   500,
		Health:      500,
		Damage:      90,
		Cooldown:    1000,
		Stamina:     100,
		StaminaCost: 90,
		Dodge:       10,
		AttackRange: 75,
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
	squadC := ecsManager.componentMap["squad"]

	entities := ecsManager.Query(ecs.BuildTag(targetC, squadC))

	for _, item := range entities {
		target := item.Components[targetC].(*TargetComponent)
		squad := item.Components[squadC].(*SquadComponent)

		targetEntities := ecsManager.Query(ecs.BuildTag(statC, squadC))
		for _, te := range targetEntities {
			targetSquad := te.Components[squadC].(*SquadComponent)

			if te.Entity.ID != item.Entity.ID && squad.Squad != targetSquad.Squad {
				target.TargetID = te.Entity.ID
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

		d2 := ((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1))
		d := float32(math.Sqrt(float64(d2)))

		if d < stats.AttackRange {
			continue
		}

		position.X += (x2 - x1) / d * dt / 30
		position.Y += (y2 - y1) / d * dt / 30
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
		d2 := ((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1))
		d := math.Sqrt(float64(d2))
		if float32(d) > stats.AttackRange {
			continue
		}

		stats.Stamina = stats.Stamina - stats.StaminaCost
		stats.Reload = 0

		if rand.Int31n(100)+1 <= targetStats.Dodge {
			continue
		}

		targetStats.Health = targetStats.Health - stats.Damage
		if targetStats.Health <= 0 {
			ecsManager.DisposeEntity(target.Entity)
		}
	}
}
