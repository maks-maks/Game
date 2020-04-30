package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

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
			q := ecsManager.Query(ecs.BuildTag(pComp))

			for _, item := range q {
				// d, _ := e.GetComponentData(ecsManager.componentMap["position"])
				data := item.Components[pComp].(*PositionComponent)

				canvas := g.GetCanvas()
				pos := g.GetCursorScreenPos()
				p1 := pos.Add(image.Pt(int(data.X), int(data.Y)))
				color := color.RGBA{200, 75, 75, 255}
				canvas.AddCircleFilled(p1, 50, color)
			}
		}),
	}
}

func loop() {
	size := g.Context.GetPlatform().DisplaySize()
	g.SingleWindow("main window", g.Layout{
		// g.Button("Show demo window", func() { demo = true }),
		g.Checkbox("Show demo window", &demo, func() {}),
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
	wnd.Main(loop)
	// renderSystem()
}
