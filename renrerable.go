package main

import (
	"github.com/bytearena/ecs"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
)

type Animatable interface {
	Update(dt float32)
	Clone() interface{}
	Reset()
}

type RenderableComponent struct {
	Node      core.INode
	Animation Animatable
}

type RenderSystem struct {
	Scene       *core.Node
	Camera      *camera.Camera
	StaticNodes []core.INode
	materials   map[string]material.IMaterial
	animations  map[string]Animatable
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		materials:  make(map[string]material.IMaterial),
		animations: make(map[string]Animatable),
	}
}

func (s *RenderSystem) Update(dt float32) {
	entities := ecsManager.Query(ecs.BuildTag(renderableC, positionC))

	for _, item := range entities {
		renderable := item.Components[renderableC].(*RenderableComponent)
		position := item.Components[positionC].(*PositionComponent)

		renderable.Node.GetNode().SetPosition(position.X/10, 0, position.Y/10)
		if renderable.Animation != nil {
			renderable.Animation.Update(dt)
		}
	}
}

func (s *RenderSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			s.startAnimation(event.DamagerID, "tank/attack")
		case *WalkStartEvent:
			s.startAnimation(event.EntityID, "tank/walk")
		case *StopEvent:
			damager := ecsManager.GetEntityByID(event.EntityID, aliveC)
			if damager == nil {
				return true
			}
			s.startAnimation(event.EntityID, "tank/idle")
		case *DeathEvent:
			s.startAnimation(event.EntityID, "tank/death")
			s.startAnimation(event.DamagerID, "tank/taunt")
		}
		return true
	})

}

func (s *RenderSystem) startAnimation(id ecs.EntityID, anim string) {
	damager := ecsManager.GetEntityByID(id, renderableC)
	if damager == nil {
		return
	}
	renderable := damager.Components[renderableC].(*RenderableComponent)

	sprite, ok := renderable.Node.GetINode().(*graphic.Sprite)
	if !ok {
		return
	}

	sprite.ClearMaterials()
	sprite.AddMaterial(sprite, s.materials[anim], 0, 0)

	renderable.Animation = s.animations[anim].Clone().(Animatable)
	renderable.Animation.Reset()
}

func (s *RenderSystem) PopulateScene() {
	entities := ecsManager.Query(ecs.BuildTag(renderableC))

	s.Scene.RemoveAll(false)
	for _, o := range s.StaticNodes {
		s.Scene.Add(o)
	}

	for _, item := range entities {
		renderable := item.Components[renderableC].(*RenderableComponent)

		s.Scene.Add(renderable.Node)
	}
}
func (s *RenderSystem) AddMaterial(key string, material material.IMaterial) {
	s.materials[key] = material
}

func (s *RenderSystem) AddAnimation(key string, a Animatable) {
	s.animations[key] = a
}
