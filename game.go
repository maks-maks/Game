package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/bytearena/ecs"
)

type Ability interface {
	Activate(id ecs.EntityID)
	Deactivate(id ecs.EntityID)
}

type DummyAbility struct{}

func (a *DummyAbility) Activate(id ecs.EntityID)   {}
func (a *DummyAbility) Deactivate(id ecs.EntityID) {}

type RageAbility struct{}

func (a *RageAbility) Activate(id ecs.EntityID) {
	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 2
	stats.Cooldown = stats.Cooldown / 3
}
func (a *RageAbility) Deactivate(id ecs.EntityID) {
	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 1
	stats.Cooldown = stats.Cooldown * 3
}

type SnipeAbility struct{}

func (a *SnipeAbility) Activate(id ecs.EntityID) {
	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 0.5
	stats.Cooldown = 3000
	stats.AttackRange = 400
	stats.Dodge = 60
	stats.Damage = 100
	stats.DodgeRange = 50
}
func (a *SnipeAbility) Deactivate(id ecs.EntityID) {
	item := ecsManager.GetEntityByID(id, statC)
	stats := item.Components[statC].(*StatsComponent)

	stats.Resist = 1
	stats.Cooldown = 400
	stats.AttackRange = 300
	stats.Dodge = 30
	stats.Damage = 20
	stats.DodgeRange = 25
}

type ReviveAbility struct{}

func (a *ReviveAbility) Activate(id ecs.EntityID) {
	item := ecsManager.GetEntityByID(id, statC, squadC)
	// stats := item.Components[statC].(*StatsComponent)
	squad := item.Components[squadC].(*SquadComponent)

	entities := ecsManager.Query(ecs.BuildTag(squadC, statC, positionC, deadC))

	for _, item := range entities {
		targetSquad := item.Components[squadC].(*SquadComponent)
		targetStats := item.Components[statC].(*StatsComponent)
		// targetPosition := item.Components[positionC].(*PositionComponent)

		if squad.Squad != targetSquad.Squad {
			continue
		}

		item.Entity.RemoveComponent(deadC)
		item.Entity.AddComponent(aliveC, &AliveComponent{})
		log = append(log, fmt.Sprintf("%d revived %d with arrow", id, item.Entity.ID))
		targetStats.Health = targetStats.MaxHealth / 2
	}

	// stats.Resist = 0.5
	// stats.Cooldown = 3000
	// stats.AttackRange = 400
	// stats.Dodge = 60
	// stats.Damage = 100
	// stats.DodgeRange = 50
}
func (a *ReviveAbility) Deactivate(id ecs.EntityID) {
	// item := ecsManager.GetEntityByID(id, statC)
	// stats := item.Components[statC].(*StatsComponent)

	// stats.Resist = 1
	// stats.Cooldown = 400
	// stats.AttackRange = 300
	// stats.Dodge = 30
	// stats.Damage = 20
	// stats.DodgeRange = 25
}

type UltimateComponent struct {
	Cooldown float32 `imgui:"%.1f ms"`
	Reload   float32 `imgui:"%.1f ms"`
	Ability  Ability
	Active   bool
	Charge   float32
	HitInc   float32
	DodgeInc float32
	HealInc  float32
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
type AliveComponent struct {
}
type DeadComponent struct {
}

type NameComponent struct {
	Name string
}
type ArrowComponent struct {
	Damage int32
}

var (
	nameC     *ecs.Component
	positionC *ecs.Component
	statC     *ecs.Component
	targetC   *ecs.Component
	squadC    *ecs.Component
	ultaC     *ecs.Component
	aliveC    *ecs.Component
	deadC     *ecs.Component
	arrowC    *ecs.Component
)

func setupECS() {
	ecsManager = NewECSManager()

	nameC = ecsManager.RegisterComponent("name", &NameComponent{})
	positionC = ecsManager.RegisterComponent("position", &PositionComponent{})
	statC = ecsManager.RegisterComponent("stats", &StatsComponent{})
	targetC = ecsManager.RegisterComponent("target", &TargetComponent{})
	squadC = ecsManager.RegisterComponent("squad", &SquadComponent{})
	ultaC = ecsManager.RegisterComponent("ulta", &UltimateComponent{})
	aliveC = ecsManager.RegisterComponent("alive", &AliveComponent{})
	deadC = ecsManager.RegisterComponent("dead", &DeadComponent{})
	arrowC = ecsManager.RegisterComponent("arrow", &ArrowComponent{})
	//createTank("Frederik", "a", 100, 100)
	//createRanger("Legolas", "b", 400, 400)
	// for i := 1; i < 5; i++ {
	createSquad("Geroi", 500, 500)
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
	ecsManager.AddComponent(e, &NameComponent{Name: n})
	ecsManager.AddComponent(e, &PositionComponent{
		X: x + rand.Float32()*200 - 100,
		Y: y + rand.Float32()*200 - 100,
	})
	ecsManager.AddComponent(e, &SquadComponent{
		Squad: squad,
	})
	ecsManager.AddComponent(e, &AliveComponent{})
	ecsManager.AddComponent(e, &UltimateComponent{
		Ability: &DummyAbility{},
	})
	ecsManager.AddComponent(e, &UltimateComponent{
		Cooldown: 5000,
		Ability:  &ReviveAbility{},
		Charge:   0,
		HealInc:  25,
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
	ecsManager.AddComponent(e, &NameComponent{Name: n})
	ecsManager.AddComponent(e, &PositionComponent{
		X: x + rand.Float32()*200 - 100,
		Y: y + rand.Float32()*200 - 100,
	})
	ecsManager.AddComponent(e, &SquadComponent{
		Squad: squad,
	})
	ecsManager.AddComponent(e, &UltimateComponent{
		Cooldown: 10000,
		Ability:  &SnipeAbility{},
		Charge:   0,
		HitInc:   5,
	})
	ecsManager.AddComponent(e, &AliveComponent{})
	//ecsManager.AddComponent(e, &UltimateComponent{
	//	Ability: &DummyAbility{},
	//})
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
	ecsManager.AddComponent(e, &NameComponent{Name: n})
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
		Charge:   0,
		HitInc:   25,
	})
	ecsManager.AddComponent(e, &AliveComponent{})
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
	entities := ecsManager.Query(ecs.BuildTag(targetC, squadC, statC, positionC, aliveC))

	for _, item := range entities {
		target := item.Components[targetC].(*TargetComponent)
		squad := item.Components[squadC].(*SquadComponent)
		stats := item.Components[statC].(*StatsComponent)
		position := item.Components[positionC].(*PositionComponent)
		targetEntities := ecsManager.Query(ecs.BuildTag(statC, squadC, positionC, aliveC))

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

type arrowSystem struct{}

func (s *arrowSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(targetC, positionC, arrowC))
	for _, item := range query {
		position := item.Components[positionC].(*PositionComponent)
		currentTarget := item.Components[targetC].(*TargetComponent)
		arrow := item.Components[arrowC].(*ArrowComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, positionC, statC)
		if target == nil {
			currentTarget.TargetID = 0
			ecsManager.DisposeEntity(item.Entity)
			continue
		}

		targetPosition := target.Components[positionC].(*PositionComponent)
		x1, y1, x2, y2 := position.X, position.Y, targetPosition.X, targetPosition.Y

		d := distance(position, targetPosition)

		position.X += (x2 - x1) / d * dt / 5
		position.Y += (y2 - y1) / d * dt / 5

		d = distance(position, targetPosition)

		if d < 15 {
			targetstats := target.Components[statC].(*StatsComponent)
			targetstats.Health -= arrow.Damage
			ecsManager.DisposeEntity(item.Entity)
			if targetstats.Health <= 0 {
				target.Entity.RemoveComponent(aliveC)
				target.Entity.AddComponent(deadC, &DeadComponent{})
				log = append(log, fmt.Sprintf("%d killed %d with arrow", item.Entity.ID, target.Entity.ID))
			}
		}
	}

}

type movementSystem struct{}

func (s *movementSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(targetC, positionC, statC, aliveC))
	for _, item := range query {
		position := item.Components[positionC].(*PositionComponent)
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statC].(*StatsComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, positionC)
		if target == nil {
			currentTarget.TargetID = 0
			continue
		}

		targetPosition := target.Components[positionC].(*PositionComponent)
		x1, y1, x2, y2 := position.X, position.Y, targetPosition.X, targetPosition.Y

		d := distance(position, targetPosition)
		position.D = d

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
	query := ecsManager.Query(ecs.BuildTag(ultaC, aliveC))

	for _, item := range query {
		ulta := item.Components[ultaC].(*UltimateComponent)
		if ulta.Reload < ulta.Cooldown {
			ulta.Reload = ulta.Reload + dt
			continue
		}
		if ulta.Charge < 100 {
			continue
		}
		ulta.Charge -= 100
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
func createArrow(targetID ecs.EntityID, x float32, y float32, damage int32) *ecs.Entity {
	e := ecsManager.NewEntity()

	ecsManager.AddComponent(e, &PositionComponent{
		X: x,
		Y: y,
	})
	ecsManager.AddComponent(e, &ArrowComponent{
		Damage: damage,
	})

	ecsManager.AddComponent(e, &TargetComponent{TargetID: targetID})
	return e
}

type battleSystem struct{}

func (s *battleSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(targetC, statC, positionC, ultaC, aliveC))

	for _, item := range query {
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statC].(*StatsComponent)
		position := item.Components[positionC].(*PositionComponent)

		target := ecsManager.GetEntityByID(currentTarget.TargetID, statC, positionC, ultaC, aliveC)
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
			ulta := item.Components[ultaC].(*UltimateComponent)
			ulta.Charge = ulta.Charge + ulta.HealInc
			if targetStats.Health > targetStats.MaxHealth {
				targetStats.Health = targetStats.MaxHealth

				continue
			}
			continue
		}

		if stats.AttackRange > 75 {
			createArrow(target.Entity.ID, position.X, position.Y, stats.Damage)
			continue
		}
		if rand.Int31n(100)+1 <= targetStats.Dodge {
			targetPosition.X += (x2 - x1) / d * 1 * targetStats.DodgeRange
			targetPosition.Y += (y2 - y1) / d * 1 * targetStats.DodgeRange
			targetUlta := target.Components[ultaC].(*UltimateComponent)
			targetUlta.Charge += targetUlta.DodgeInc
			continue
		}

		targetStats.Health = targetStats.Health - int32((float32(stats.Damage) / targetStats.Resist))
		ulta := item.Components[ultaC].(*UltimateComponent)
		ulta.Charge = ulta.Charge + ulta.HitInc
		if targetStats.Health <= 0 {
			target.Entity.RemoveComponent(aliveC)
			target.Entity.AddComponent(deadC, &DeadComponent{})
			log = append(log, fmt.Sprintf("%d killed %d", item.Entity.ID, target.Entity.ID))
		}
	}
}
