package main

import (
	"github.com/bytearena/ecs"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

type Animatable interface {
	Update(dt float32)
	Clone() interface{}
	Reset()
}

type RenderableComponent struct {
	Node      core.INode
	Animation Animatable
	Base      string
}

type RenderSystem struct {
	Scene       *core.Node
	Camera      *camera.Camera
	StaticNodes []core.INode
	animations  map[string]Animatable
	textures    map[string]*texture.Texture2D
	defaultMat  material.IMaterial
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		animations: make(map[string]Animatable),
		textures:   make(map[string]*texture.Texture2D),
		defaultMat: material.NewStandard(&math32.Color{1, 1, 1}),
	}
}

func (s *RenderSystem) Update(dt float32) {
	entities := ecsManager.Query(ecs.BuildTag(renderableC, positionC))

	for _, item := range entities {
		renderable := item.Components[renderableC].(*RenderableComponent)
		position := item.Components[positionC].(*PositionComponent)

		renderable.Node.GetNode().SetPosition(position.X/10, 0, position.Y/10)

		if sprite, ok := renderable.Node.GetINode().(*graphic.Sprite); ok {
			for _, gm := range sprite.Materials() {
				if updater, ok2 := gm.IMaterial().(interface {
					Update(float32)
				}); ok2 {
					updater.Update(dt)
				}
			}
		}
	}
}

func (s *RenderSystem) ProcessEvents(b EventBus) {
	b.Iterate(func(e Event) bool {
		switch event := e.(type) {
		case *HitEvent:
			s.startAnimation(event.DamagerID, "attack")
		case *WalkStartEvent:
			s.startAnimation(event.EntityID, "walk")
		case *StopEvent:
			damager := ecsManager.GetEntityByID(event.EntityID, aliveC)
			if damager == nil {
				return true
			}
			s.startAnimation(event.EntityID, "idle")
		case *DeathEvent:
			s.startAnimation(event.EntityID, "death")
			s.startAnimation(event.DamagerID, "taunt")
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

	anim = renderable.Base + "/" + anim

	renderable.Animation = s.animations[anim].Clone().(Animatable)
	renderable.Animation.Reset()

	sprite, ok := renderable.Node.GetINode().(*graphic.Sprite)
	if !ok {
		return
	}

	sprite.ClearMaterials()

	mat := material.NewStandard(&math32.Color{1, 1, 1})
	mat.AddTexture(s.textures[anim])
	mat.SetOpacity(1)
	mat.SetTransparent(true)
	sprite.AddMaterial(sprite, &AnimatedMaterial{
		IMaterial: mat,
		Animation: renderable.Animation.(UpdaterRenderSetuper),
	}, 0, 0)
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

func (s *RenderSystem) AddAnimation(key string, a Animatable) {
	s.animations[key] = a
}

func (s *RenderSystem) AddTexture(key string, texture *texture.Texture2D) {
	s.textures[key] = texture
}
