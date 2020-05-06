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
			pComp := ecsManager.componentMap["position"]
			sComp := ecsManager.componentMap["stats"]
			q := ecsManager.Query(ecs.BuildTag(pComp))

			for _, item := range q {
				// d, _ := e.GetComponentData(ecsManager.componentMap["position"])
				data := item.Components[pComp].(*PositionComponent)

				canvas := g.GetCanvas()
				pos := g.GetCursorScreenPos()
				p0 := pos.Add(image.Pt(int(data.X), int(data.Y)))
				circleColor := color.RGBA{255, 255, 255, 255}
				r := 25
				canvas.AddCircleFilled(p0, float32(r), circleColor)
				canvas.AddRectFilled(
					p0.Add(image.Pt(-3, 20)), // top left
					p0.Add(image.Pt(3, 150)), // bottom right
					circleColor,              // color
					0,
					giu.CornerFlags_All,
				)
				// Healthbar
				s, ok := item.Entity.GetComponentData(sComp)
				if ok {
					stats := s.(*StatsComponent)
					width := r * 2 * int(stats.Health) / int(stats.MaxHealth)

					pMin := p0.Add(image.Pt(-r, -r-5))
					pMax := pMin.Add(image.Pt(int(width), 10))
					hbColor := color.RGBA{200, 0, 0, 255}
					canvas.AddRectFilled(pMin, pMax, hbColor, 0, giu.CornerFlags_All)

					textPos := pMin
					textColor := color.RGBA{0, 0, 0, 255}
					canvas.AddText(textPos, textColor, fmt.Sprintf("%d", stats.Health))
				}
			}
		}),
	}
}

func loop() {
	size := g.Context.GetPlatform().DisplaySize()
	g.SingleWindow("main window", g.Layout{
		// g.Button("Show demo window", func() { demo = true }),
		g.Line(
			g.Checkbox("Show demo window", &demo, func() {}),
			g.SliderInt("Delay", &delayMs, 10, 500, "%d ms"),
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
		&battleSystem{},
	}

	go func() {
		t := time.Now()
		for {
			dt := float32(time.Since(t).Milliseconds())
			t = time.Now()

			for _, s := range systems {
				s.Update(dt)
			}

			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			g.Update()
		}
	}()

	wnd.Main(loop)
	// renderSystem()
}
