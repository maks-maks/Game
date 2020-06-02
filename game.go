package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/bytearena/ecs"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
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
	squad := item.Components[squadC].(*SquadComponent)

	entities := ecsManager.Query(ecs.BuildTag(squadC, statC, positionC, deadC))

	for _, item := range entities {
		targetSquad := item.Components[squadC].(*SquadComponent)

		if squad.Squad != targetSquad.Squad {
			continue
		}

		ecsManager.events.Schedule(&ReviveEvent{
			EntityID: item.Entity.ID,
		}, 0)
	}
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
	ArrowDodge  int32
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
	Damage    int32
	DamagerID ecs.EntityID
}

var (
	nameC       *ecs.Component
	positionC   *ecs.Component
	statC       *ecs.Component
	targetC     *ecs.Component
	squadC      *ecs.Component
	ultaC       *ecs.Component
	aliveC      *ecs.Component
	deadC       *ecs.Component
	arrowC      *ecs.Component
	stateC      *ecs.Component
	renderableC *ecs.Component
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
	stateC = ecsManager.RegisterComponent("state", &StateComponent{})
	renderableC = ecsManager.RegisterComponent("renderable", &RenderableComponent{})
	//createTank("Frederik", "a", 100, 100)
	//createRanger("Legolas", "b", 400, 400)
	// for i := 1; i < 5; i++ {
	createSquad("Geroi", 200, 200)
	createSquad("Sandali", -200, -200)
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

	createRanger("Legolas9", squad, x, y)

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
	ecsManager.AddComponent(e, &StateComponent{
		State: "idle",
	})
	geom := geometry.NewSphere(1, 24, 24)
	mat := material.NewPhysical()
	mat.SetBaseColorFactor(math32.NewColor4("Yellow", 1))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetRotation(-math32.Pi/2, 0, 0)
	ecsManager.AddComponent(e, &RenderableComponent{
		Node: mesh,
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
		ArrowDodge:  10,
	})
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}
func createRanger(n string, squad string, x float32, y float32) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &NameComponent{Name: n})
	ecsManager.AddComponent(e, &StateComponent{
		State: "idle",
	})
	geom := geometry.NewSphere(1, 24, 24)
	mat := material.NewPhysical()
	mat.SetBaseColorFactor(math32.NewColor4("Yellow", 1))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetRotation(-math32.Pi/2, 0, 0)
	ecsManager.AddComponent(e, &RenderableComponent{
		Node: mesh,
	})
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
		ArrowDodge:  20,
	})
	ecsManager.AddComponent(e, &TargetComponent{})
	return e
}

func createTank(n string, squad string, x float32, y float32) *ecs.Entity {
	e := ecsManager.NewEntity()
	ecsManager.AddComponent(e, &NameComponent{Name: n})
	ecsManager.AddComponent(e, &StateComponent{
		State: "idle",
	})
	geom := geometry.NewSphere(1, 24, 24)
	mat := material.NewPhysical()
	mat.SetBaseColorFactor(math32.NewColor4("Yellow", 1))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetRotation(-math32.Pi/2, 0, 0)
	ecsManager.AddComponent(e, &RenderableComponent{
		Node: mesh,
	})
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
		ArrowDodge:  5,
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

type chasingSystem struct{}

func (s *chasingSystem) Update(dt float32) {
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
		if d < stats.AttackRange {
			position.XSpeed = 0
			position.YSpeed = 0
			continue
		}

		position.XSpeed = (x2 - x1) / d / 30
		position.YSpeed = (y2 - y1) / d / 30
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

		if d < 15 {
			ecsManager.DisposeEntity(item.Entity)

			targetstats := target.Components[statC].(*StatsComponent)
			if targetstats.Health <= 0 {
				continue
			}
			if rand.Int31n(100)+1 <= targetstats.ArrowDodge {
				dx := (x2 - x1) / d * 1 * targetstats.DodgeRange
				dy := (y2 - y1) / d * 1 * targetstats.DodgeRange

				targetPosition.X += -dy
				targetPosition.Y += dx

				continue
			}

			ecsManager.events.Schedule(&HitEvent{
				DamagerID: arrow.DamagerID,
				TargetID:  target.Entity.ID,
				Damage:    arrow.Damage,
			}, 0)
		}
	}

}

type movementSystem struct{}

func (s *movementSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(positionC))
	for _, item := range query {
		position := item.Components[positionC].(*PositionComponent)
		position.X += dt * position.XSpeed
		position.Y += dt * position.YSpeed
	}
}

func distance(a, b *PositionComponent) float32 {
	d2 := ((b.X - a.X) * (b.X - a.X)) + ((b.Y - a.Y) * (b.Y - a.Y))
	return float32(math.Sqrt(float64(d2)))
}

func distanceXY(x1, y1, x2, y2 float32) float32 {
	d2 := (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)
	return float32(math.Sqrt(float64(d2)))
}

type ultimatesSystem struct{}

func (s *ultimatesSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *DamageEvent:
			entity := ecsManager.GetEntityByID(event.TargetID, ultaC, aliveC)
			if entity == nil {
				return true
			}
			ulta, ok := entity.Components[ultaC].(*UltimateComponent)
			if !ok {
				return true
			}
			ulta.Charge += ulta.HitInc
		case *DodgeEvent:
			entity := ecsManager.GetEntityByID(event.EntityID, ultaC, aliveC)
			ulta, ok := entity.Components[ultaC].(*UltimateComponent)
			if !ok {
				return true
			}
			ulta.Charge += ulta.DodgeInc
		case *HealEvent:
			entity := ecsManager.GetEntityByID(event.TargetID, ultaC, aliveC)
			ulta, ok := entity.Components[ultaC].(*UltimateComponent)
			if !ok {
				return true
			}

			ulta.Charge += ulta.HealInc
		}
		return true
	})
}

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

func createArrow(damagerID, targetID ecs.EntityID, x float32, y float32, damage int32, x2 float32, y2 float32) *ecs.Entity {
	d := distanceXY(x, y, x2, y2)

	e := ecsManager.NewEntity()

	ecsManager.AddComponent(e, &PositionComponent{
		X:      x,
		Y:      y,
		XSpeed: (x2 - x) / d,
		YSpeed: (y2 - y) / d,
	})
	ecsManager.AddComponent(e, &ArrowComponent{
		Damage:    damage,
		DamagerID: damagerID,
	})

	ecsManager.AddComponent(e, &TargetComponent{TargetID: targetID})

	geom := geometry.NewSphere(0.5, 24, 24)
	mat := material.NewPhysical()
	mat.SetBaseColorFactor(math32.NewColor4("Red", 1))
	mesh := graphic.NewMesh(geom, mat)
	ecsManager.AddComponent(e, &RenderableComponent{
		Node: mesh,
	})

	return e
}

type dodgeSystem struct{}

func (s *dodgeSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *DodgeEvent:
			entity := ecsManager.GetEntityByID(event.EntityID, statC, positionC, aliveC)
			position := entity.Components[positionC].(*PositionComponent)
			stats := entity.Components[statC].(*StatsComponent)

			position.X += event.Direction[0] * stats.DodgeRange
			position.Y += event.Direction[1] * stats.DodgeRange
		}
		return true
	})
}

func (s *dodgeSystem) Update(dt float32) {}

type aliveSystem struct{}

func (s *aliveSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			target := ecsManager.GetEntityByID(event.TargetID, statC, aliveC, stateC)
			if target == nil {
				return true
			}
			targetStats := target.Components[statC].(*StatsComponent)
			if targetStats.Health <= 0 {
				target.Entity.RemoveComponent(aliveC)
				target.Entity.AddComponent(deadC, &DeadComponent{})

				ecsManager.events.Schedule(&DeathEvent{
					EntityID: target.Entity.ID,
				}, 0)
			}
		case *ReviveEvent:
			target := ecsManager.GetEntityByID(event.EntityID, statC, deadC, stateC)
			target.Entity.RemoveComponent(deadC)
			target.Entity.AddComponent(aliveC, &AliveComponent{})

			targetStats := target.Components[statC].(*StatsComponent)
			targetStats.Health = targetStats.MaxHealth / 2
		}
		return true
	})
}

func (s *aliveSystem) Update(dt float32) {}

type battleSystem struct{}

func (s *battleSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			target := ecsManager.GetEntityByID(event.TargetID, statC, positionC, aliveC, stateC)
			if target == nil {
				return true
			}
			targetStats := target.Components[statC].(*StatsComponent)

			if rand.Int31n(100)+1 <= targetStats.Dodge {
				targetPosition := target.Components[positionC].(*PositionComponent)
				x1, y1, x2, y2 := event.DamagerPosition[0], event.DamagerPosition[1], targetPosition.X, targetPosition.Y
				d := distanceXY(x1, y1, x2, y2)
				var dir [2]float32
				// TODO tank, ranger, arrow different dodge direction
				dir[0] += (x2 - x1) / d
				dir[1] += (y2 - y1) / d
				ecsManager.events.Schedule(&DodgeEvent{
					DamagerID: event.DamagerID,
					EntityID:  target.Entity.ID,
					Direction: dir,
				}, 0)

				return true
			}

			targetStats.Health = targetStats.Health - int32((float32(event.Damage) / targetStats.Resist))
			ecsManager.events.Schedule(&DamageEvent{
				DamagerID: event.DamagerID,
				TargetID:  target.Entity.ID,
				Damage:    event.Damage,
			}, 0)

			targetState := target.Components[stateC].(*StateComponent)
			targetState.State = "idle"
		case *HealEvent:
			healer := ecsManager.GetEntityByID(event.HealerID, statC, positionC, aliveC, stateC)
			if healer == nil {
				return true
			}
			healerStats := healer.Components[statC].(*StatsComponent)

			target := ecsManager.GetEntityByID(event.TargetID, statC, positionC, aliveC, stateC)
			if target == nil {
				return true
			}
			targetStats := target.Components[statC].(*StatsComponent)

			targetStats.Health = targetStats.Health + healerStats.Heal
			if targetStats.Health > targetStats.MaxHealth {
				targetStats.Health = targetStats.MaxHealth
			}
		}
		return true
	})
}

func (s *battleSystem) Update(dt float32) {
	query := ecsManager.Query(ecs.BuildTag(targetC, statC, positionC, ultaC, aliveC, stateC))

	for _, item := range query {
		currentTarget := item.Components[targetC].(*TargetComponent)
		stats := item.Components[statC].(*StatsComponent)
		position := item.Components[positionC].(*PositionComponent)
		state := item.Components[stateC].(*StateComponent)

		if state.State == "attack" {
			continue
		}

		target := ecsManager.GetEntityByID(currentTarget.TargetID, statC, positionC, aliveC)
		if target == nil {
			currentTarget.TargetID = 0
			continue
		}

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
		d := distance(position, targetPosition)
		if d > stats.AttackRange {
			continue
		}

		stats.Stamina = stats.Stamina - stats.StaminaCost
		stats.Reload = 0
		if stats.Heal > 0 {
			ecsManager.events.Schedule(&HealEvent{
				HealerID: item.Entity.ID,
				TargetID: target.Entity.ID,
				Health:   stats.Heal,
			}, 0)
			continue
		}

		if stats.AttackRange > 75 {
			createArrow(item.Entity.ID, target.Entity.ID, position.X, position.Y, stats.Damage, targetPosition.X, targetPosition.Y)
			continue
		}

		state.State = "attack"
		ecsManager.events.Schedule(&HitEvent{
			DamagerID:       item.Entity.ID,
			TargetID:        target.Entity.ID,
			Damage:          stats.Damage,
			DamagerPosition: [2]float32{position.X, position.Y},
		}, 100)
	}
}

type HitEvent struct {
	DamagerID       ecs.EntityID
	TargetID        ecs.EntityID
	Damage          int32
	DamagerPosition [2]float32
}

func (e HitEvent) String() string {
	return fmt.Sprintf("{HitEvent %d from %d to %d pos=%v}", e.Damage, e.DamagerID, e.TargetID, e.DamagerPosition)
}

type DamageEvent struct {
	DamagerID ecs.EntityID
	TargetID  ecs.EntityID
	Damage    int32
}

func (e DamageEvent) String() string {
	return fmt.Sprintf("{DamageEvent %d from %d to %d}", e.Damage, e.DamagerID, e.TargetID)
}

type DodgeEvent struct {
	DamagerID ecs.EntityID
	EntityID  ecs.EntityID
	Direction [2]float32
}

func (e DodgeEvent) String() string {
	return fmt.Sprintf("{DodgeEvent %d from %d dir=%v}", e.EntityID, e.DamagerID, e.Direction)
}

type HealEvent struct {
	HealerID ecs.EntityID
	TargetID ecs.EntityID
	Health   int32
}

func (e HealEvent) String() string {
	return fmt.Sprintf("{HealEvent %d from %d to %d}", e.Health, e.HealerID, e.TargetID)
}

type DeathEvent struct {
	EntityID ecs.EntityID
}

func (e DeathEvent) String() string {
	return fmt.Sprintf("{DeathEvent %d}", e.EntityID)
}

type ReviveEvent struct {
	EntityID ecs.EntityID
}

func (e ReviveEvent) String() string {
	return fmt.Sprintf("{ReviveEvent %d}", e.EntityID)
}
