package main

import (
	"github.com/bytearena/ecs"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
)

type RenderableComponent struct {
	Node core.INode
}

type renderableSystem struct {
	Scene       *core.Node
	Camera      *camera.Camera
	StaticNodes []core.INode
}

func (s *renderableSystem) Update(dt float32) {
	entities := ecsManager.Query(ecs.BuildTag(renderableC, positionC))

	for _, item := range entities {
		renderable := item.Components[renderableC].(*RenderableComponent)
		position := item.Components[positionC].(*PositionComponent)

		renderable.Node.GetNode().SetPosition(position.X/10, 0, position.Y/10)
	}
}

func (s *renderableSystem) PopulateScene() {
	entities := ecsManager.Query(ecs.BuildTag(renderableC))

	s.Scene.RemoveAll(true)
	for _, o := range s.StaticNodes {
		s.Scene.Add(o)
	}

	for _, item := range entities {
		renderable := item.Components[renderableC].(*RenderableComponent)

		s.Scene.Add(renderable.Node)
	}
}
