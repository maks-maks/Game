package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/bytearena/ecs"
)

type PositionComponent struct {
	X float32
	Y float32
}

func renderSystem() {
	entities := ecsManager.Query(0).Entities()

	for _, v := range entities {
		fmt.Println(v)
	}
}

func gameCanvas() *g.Layout {
	return &g.Layout{
		// g.Label("Canvas demo"),
		g.Custom(func() {
			q := ecsManager.Query(ecs.BuildTag(positionC, statC))

			for _, item := range q {
				data := item.Components[positionC].(*PositionComponent)
				stats := item.Components[statC].(*StatsComponent)

				canvas := g.GetCanvas()
				pos := g.GetCursorScreenPos()
				p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))

				opacity := uint8(255)
				if stats.Health <= 0 {
					opacity = 25
				}

				var circleColor color.RGBA
				if stats.Heal > 0 {
					circleColor = color.RGBA{75, 120, 210, opacity}
				} else if stats.AttackRange > 100 {
					circleColor = color.RGBA{200, 0, 100, opacity}
				} else {
					circleColor = color.RGBA{255, 255, 255, opacity}
				}
				r := 25
				canvas.AddCircleFilled(p0, float32(r), circleColor)
				if curEntityID == item.Entity.ID {
					canvas.AddCircle(p0, float32(r+3), color.RGBA{50, 255, 50, 255}, 3)
				}
				// Healthbar
				if stats.Health > 0 {
					width := r * 2 * int(stats.Health) / int(stats.MaxHealth)

					pMin := p0.Add(image.Pt(-r, -r-5))
					pMax := pMin.Add(image.Pt(int(width), 10))
					hbColor := color.RGBA{200, 0, 0, opacity}
					canvas.AddRectFilled(pMin, pMax, hbColor, 0, giu.CornerFlags_All)

					textPos := pMin.Add(image.Pt(0, -2))
					textColor := color.RGBA{0, 0, 0, opacity}
					canvas.AddText(textPos, textColor, fmt.Sprintf("%d", stats.Health))
				}
			}

			q = ecsManager.Query(ecs.BuildTag(positionC, arrowC))

			for _, item := range q {
				data := item.Components[positionC].(*PositionComponent)

				canvas := g.GetCanvas()
				pos := g.GetCursorScreenPos()
				p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))

				opacity := uint8(255)

				var circleColor color.RGBA
				circleColor = color.RGBA{255, 255, 255, opacity}
				r := 5
				canvas.AddCircleFilled(p0, float32(r), circleColor)
			}
		}),
	}
}

var speedMultiplier float32 = 1
var paused float32 = 0

func loop() {
	size := g.Context.GetPlatform().DisplaySize()
	var playPauseButton g.Widget
	if paused != 0 {
		playPauseButton = g.Button(" pause ", func() { paused = 0 })
	} else {
		playPauseButton = g.Button(" play  ", func() { paused = 1 })
	}
	g.SingleWindow("main window", g.Layout{
		// g.Button("Show demo window", func() { demo = true }),
		g.Line(
			g.Checkbox("Show demo window", &demo, func() {}),
			// g.SliderInt("Delay", &delayMs, 10, 500, "%d ms"),
			playPauseButton,
			g.RadioButton("0.2x", speedMultiplier == 0.2, func() { speedMultiplier = 0.2 }),
			g.RadioButton("1x", speedMultiplier == 1, func() { speedMultiplier = 1 }),
			g.RadioButton("10x", speedMultiplier == 10, func() { speedMultiplier = 10 }),
		),
		g.SplitLayout("Split", g.DirectionHorizontal, false, 300,
			leftPanel(),
			g.SplitLayout("Split2", g.DirectionHorizontal, false, size[0]-600,
				gameCanvas(),
				rightPanel(),
			),
		),
	})
	if demo {
		imgui.ShowDemoWindow(&demo)
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	setupECS()
	wnd := g.NewMasterWindow("App", 1000, 500, 0, nil)

	systems := []System{
		&targetingSystem{},
		&movementSystem{},
		&ultimatesSystem{},
		&battleSystem{},
		&arrowSystem{},
	}

	go func() {
		t := time.Now()
		for {
			dt := float32(time.Since(t).Milliseconds())
			t = time.Now()

			for _, s := range systems {
				s.Update(dt * speedMultiplier * (1 - paused))
			}

			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			g.Update()
		}
	}()

	wnd.Main(loop)
	// renderSystem()
}
