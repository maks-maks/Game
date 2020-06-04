package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/maks-maks/Game/imguipod"
)

var log []string

type PositionComponent struct {
	X      float32
	Y      float32
	XSpeed float32
	YSpeed float32
	Debug  string
}

func renderSystem() {
	entities := ecsManager.Query(0).Entities()

	for _, v := range entities {
		fmt.Println(v)
	}
}

// func gameCanvas() *g.Layout {
// 	return &g.Layout{
// 		// g.Label("Canvas demo"),
// 		g.Custom(func() {
// 			q := ecsManager.Query(ecs.BuildTag(positionC, statC, stateC))

// 			for _, item := range q {
// 				data := item.Components[positionC].(*PositionComponent)
// 				stats := item.Components[statC].(*StatsComponent)

// 				canvas := g.GetCanvas()
// 				pos := g.GetCursorScreenPos()
// 				p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))

// 				opacity := uint8(255)
// 				if stats.Health <= 0 {
// 					opacity = 25
// 				}

// 				var circleColor color.RGBA
// 				if stats.Heal > 0 {
// 					circleColor = color.RGBA{75, 120, 210, opacity}
// 				} else if stats.AttackRange > 100 {
// 					circleColor = color.RGBA{200, 0, 100, opacity}
// 				} else {
// 					circleColor = color.RGBA{255, 255, 255, opacity}
// 				}
// 				r := 25
// 				canvas.AddCircleFilled(p0, float32(r), circleColor)
// 				if curEntityID == item.Entity.ID {
// 					canvas.AddCircle(p0, float32(r+3), color.RGBA{50, 255, 50, 255}, 3)
// 				}
// 				// Healthbar
// 				if stats.Health > 0 {
// 					width := r * 2 * int(stats.Health) / int(stats.MaxHealth)

// 					pMin := p0.Add(image.Pt(-r, -r-5))
// 					pMax := pMin.Add(image.Pt(int(width), 10))
// 					hbColor := color.RGBA{200, 0, 0, opacity}
// 					canvas.AddRectFilled(pMin, pMax, hbColor, 0, giu.CornerFlags_All)

// 					textPos := pMin.Add(image.Pt(0, -2))
// 					textColor := color.RGBA{0, 0, 0, opacity}
// 					canvas.AddText(textPos, textColor, fmt.Sprintf("%d", stats.Health))
// 				}

// 				if true {
// 					state := item.Components[stateC].(*StateComponent)
// 					statePos := p0.Add(image.Pt(-r, +r+5))
// 					canvas.AddText(statePos, color.RGBA{0, 0, 0, 255}, fmt.Sprintf("%v", state.State))
// 				}
// 			}

// 			q = ecsManager.Query(ecs.BuildTag(positionC, arrowC))

// 			for _, item := range q {
// 				data := item.Components[positionC].(*PositionComponent)

// 				canvas := g.GetCanvas()
// 				pos := g.GetCursorScreenPos()
// 				p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))

// 				opacity := uint8(255)

// 				var circleColor color.RGBA
// 				circleColor = color.RGBA{255, 255, 255, opacity}
// 				// r := 5
// 				n := float32(20)
// 				vd := distanceXY(0, 0, data.XSpeed, data.YSpeed)
// 				vx := data.XSpeed / vd
// 				vy := data.YSpeed / vd
// 				data.Debug = fmt.Sprintf("v(%v %v) d:%v  vx:%v vy:%v", data.XSpeed, data.YSpeed, vd, vx, vy)
// 				pb := p0.Add(image.Pt(int(-vx*n), int(-vy*n)))

// 				p1 := pb.Add(image.Pt(int(-vy*n/2), int(vx*n/2)))
// 				p2 := pb.Add(image.Pt(int(vy*n/2), int(-vx*n/2)))
// 				p3 := p0.Add(image.Pt(int(-vx*n*2), int(-vy*n*2)))

// 				canvas.AddTriangleFilled(p0, p1, p2, circleColor)
// 				canvas.AddLine(p0, p3, circleColor, 5)
// 			}

// 			if showSpeed {
// 				q = ecsManager.Query(ecs.BuildTag(positionC))

// 				for _, item := range q {
// 					data := item.Components[positionC].(*PositionComponent)

// 					canvas := g.GetCanvas()
// 					pos := g.GetCursorScreenPos()
// 					p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))

// 					color := color.RGBA{255, 100, 100, 255}
// 					// canvas.AddCircleFilled(p0, float32(r), circleColor)
// 					canvas.AddLine(p0, p0.Add(image.Pt(int(data.XSpeed*300), int(data.YSpeed*300))), color, 2)
// 				}
// 			}
// 		}),
// 	}
// }

var speedMultiplier float32 = 1
var paused float32 = 0
var showSpeed bool = false

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	a := NewApplication(1280, 720, "game")

	scene := core.NewNode()

	cam := camera.New(1)
	cam.SetPosition(0, 100, 0)
	cam.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 0, -1})

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 1.0)

	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 10.0)
	pointLight1.SetPosition(1, 20, 50)

	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 10.0)
	pointLight2.SetPosition(1, -20, -50)

	setupECS()
	renderSystem := newRenderableSystem()
	renderSystem.Scene = scene
	renderSystem.Camera = cam
	renderSystem.StaticNodes = []core.INode{
		cam,
		ambLight,
		pointLight1,
		pointLight2,
		helper.NewAxes(1),
	}
	animationSystem := newAnimationSystem()
	systems := []System{
		&targetingSystem{},
		&chasingSystem{},
		&movementSystem{},
		&ultimatesSystem{},
		&battleSystem{},
		&dodgeSystem{},
		&arrowSystem{},
		&aliveSystem{},
		animationSystem,
		renderSystem,
	}

	loadSprites(animationSystem, renderSystem)

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 1.0, 1.0)

	glfwWindow := a.IWindow.(*window.GlfwWindow).Window
	pod := imguipod.New(glfwWindow)

	gui := NewGUI(pod)

	t := time.Now()

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		dt := float32(time.Since(t).Milliseconds())
		t = time.Now()

		ecsManager.events.AdvanceScheduled(t)

		for _, s := range systems {
			if ep, ok := s.(interface {
				ProcessEvents(EventBus)
			}); ok {
				ep.ProcessEvents(ecsManager.events)
			}
		}

		for _, s := range systems {
			s.Update(dt * speedMultiplier * (1 - paused))
		}

		ecsManager.events.ClearQueue()

		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderSystem.PopulateScene()
		renderer.Render(scene, cam)

		pod.BeginImgui()
		gui.Render()
		pod.EndImgui()
	})
}

func loadSprites(as *AnimationSystem, rs *renderableSystem) {
	loadSprite("spr_Idle_strip.png", "tank/idle", 16, as, rs)
	loadSprite("spr_Attack_strip.png", "tank/attack", 30, as, rs)
}

func loadSprite(fileName string, state string, frames int, as *AnimationSystem, rs *renderableSystem) {
	tex1, err := texture.NewTexture2DFromImage(fileName)
	if err != nil {
		panic(err)
	}
	tex1.SetMagFilter(gls.NEAREST)
	anim1 := texture.NewAnimator(tex1, frames, 1)
	anim1.SetDispTime(16666 * time.Microsecond)
	as.AddAnimation(state, anim1)

	mat1 := material.NewStandard(&math32.Color{1, 1, 1})
	mat1.AddTexture(tex1)
	mat1.SetOpacity(1)
	mat1.SetTransparent(true)
	// s1 := graphic.NewSprite(17, 9.6, mat1)
	rs.AddMaterial(state, mat1)
}
