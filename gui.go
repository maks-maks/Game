package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bytearena/ecs"
	"github.com/inkyblackness/imgui-go/v2"
	"github.com/maks-maks/Game/imguipod"
)

type Widget interface {
	Build()
}

type BaseWidget struct {
	pod *imguipod.ImguiPod
}

type GUI struct {
	pod     *imguipod.ImguiPod
	widgets []Widget
}

func NewGUI(pod *imguipod.ImguiPod) *GUI {
	bw := BaseWidget{pod}
	return &GUI{
		pod: pod,
		widgets: []Widget{
			&topToolBar{bw},
			&sceneInspectorWidget{bw},
			&componentInspector{bw},
		},
	}
}

func (g *GUI) Render() {
	for _, w := range g.widgets {
		w.Build()
	}
}

type sceneInspectorWidget struct {
	BaseWidget
}

func (w sceneInspectorWidget) Build() {
	ss := w.pod.Platform.DisplaySize()
	wh := imgui.CalcTextSize("T", false, 1000)
	top := wh.Y + imgui.CurrentStyle().FramePadding().Y*2
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: top})
	imgui.SetNextWindowSize(imgui.Vec2{X: 200, Y: ss[1] - top})

	imgui.BeginV("SceneInspector", nil, imgui.WindowFlagsNoMove)

	if imgui.BeginTabBar("LeftTabBar") {

		if imgui.BeginTabItem("Entities") {
			var tag ecs.Tag = 0
			q := ecsManager.Query(tag)
			items := make([]string, 0, len(q))

			for _, v := range q {
				nameQ, hasName := v.Entity.GetComponentData(nameC)
				squadQ, hasSquad := v.Entity.GetComponentData(squadC)
				if hasName && hasSquad {
					name := nameQ.(*NameComponent)
					squad := squadQ.(*SquadComponent)
					items = append(items, fmt.Sprintf("%d - %s - %s", v.Entity.ID, name.Name, squad.Squad))
				} else if hasName {
					name := nameQ.(*NameComponent)
					items = append(items, fmt.Sprintf("%d - %s", v.Entity.ID, name.Name))
				} else {
					items = append(items, fmt.Sprintf("%d", v.Entity.ID))
				}
			}
			ListBoxF("Entities", len(items), func(i int) {
				selected := i == curEntityI
				if imgui.SelectableV(items[i], selected, 0, imgui.Vec2{}) {
					curEntityI = i
					curEntityID = q.Entities()[curEntityI].ID
				}
			})
			imgui.EndTabItem()
		}

		if imgui.BeginTabItem("Log") {
			ListBoxF("Entities", len(log), func(i int) {
				imgui.Text(log[i])
			})

			imgui.EndTabItem()
		}

		imgui.EndTabBar()
	}

	imgui.End()
}

func ListBoxF(id string, length int, f func(i int)) bool {
	imgui.BeginChildV(id, imgui.Vec2{}, true, 0)

	var clipper imgui.ListClipper
	clipper.Begin(length)

	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			f(i)
		}
	}

	clipper.End()

	imgui.EndChild()
	return false
}

var curEntityI int
var curEntityID ecs.EntityID = 0
var delayMs int32 = 30
var demo = false

type topToolBar struct {
	BaseWidget
}

func (w topToolBar) Build() {
	ss := w.pod.Platform.DisplaySize()
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSize(imgui.Vec2{X: ss[0], Y: 0})
	open := true

	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0)
	imgui.BeginV("topbar", &open, imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoMove)

	imgui.Checkbox("Show demo window", &demo)

	imgui.SameLineV(0, 10)
	if paused == 0 {
		if imgui.Button("play") {
			paused = 1
		}
	} else {
		if imgui.Button("pause") {
			paused = 0
		}
	}

	imgui.SameLineV(0, 10)
	if imgui.RadioButton("0.2x", speedMultiplier == 0.2) {
		speedMultiplier = 0.2
	}
	imgui.SameLine()
	if imgui.RadioButton("1x", speedMultiplier == 1) {
		speedMultiplier = 1
	}
	imgui.SameLine()
	if imgui.RadioButton("10x", speedMultiplier == 10) {
		speedMultiplier = 10
	}

	imgui.SameLineV(0, 10)
	imgui.Checkbox("Show speed", &showSpeed)

	imgui.End()
	imgui.PopStyleVar()

	if demo {
		imgui.ShowDemoWindow(&demo)
	}
}

type componentInspector struct {
	BaseWidget
}

func (w componentInspector) Build() {
	ss := w.pod.Platform.DisplaySize()
	wh := imgui.CalcTextSize("T", false, 1000)
	top := wh.Y + imgui.CurrentStyle().FramePadding().Y*2
	imgui.SetNextWindowPos(imgui.Vec2{X: ss[0] - 200, Y: top})
	imgui.SetNextWindowSize(imgui.Vec2{X: 200, Y: ss[1] - top})

	imgui.BeginV("ComponentInspector", nil, imgui.WindowFlagsNoMove)
	defer imgui.End()

	q := ecsManager.GetEntityByID(curEntityID)
	if q == nil {
		return
	}

	for _, component := range ecsManager.components {
		if q.Entity.HasComponent(component) {
			if imgui.TreeNodeV(ecsManager.ComponentName(component), imgui.TreeNodeFlagsDefaultOpen) {
				buildComponentWidget(q.Entity, component)
				imgui.TreePop()
			}
		}
	}
}

func buildComponentWidget(e *ecs.Entity, c *ecs.Component) {
	d, _ := e.GetComponentData(c)
	val := reflect.ValueOf(d)
	val = reflect.Indirect(val)

	buildStructWidget(val)
}

func buildStructWidget(val reflect.Value) {
	if !val.IsValid() {
		return
	}
	val = reflect.Indirect(val)
	typ := val.Type()

	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			ft := f.Type
			vf := val.Field(i)
			kind := ft.Kind()

			tag, ok := f.Tag.Lookup("imgui")
			format := ""
			if ok {
				format = strings.Split(tag, ",")[0]
			}

			switch kind {
			case reflect.String:
				imgui.InputText(f.Name, vf.Addr().Interface().(*string))
			case reflect.Float32:
				format = stringOrDefault(format, "%.3f")
				imgui.DragFloatV(f.Name, vf.Addr().Interface().(*float32), 1, 0, 0, format, 1)
			case reflect.Int32:
				format = stringOrDefault(format, "%d")
				imgui.DragIntV(f.Name, vf.Addr().Interface().(*int32), 1.0, 0, 0, format)
			case reflect.Uint32:
				format = stringOrDefault(format, "%d")
				imgui.LabelText(f.Name, fmt.Sprintf("%d", vf.Interface().(ecs.EntityID)))
			case reflect.Bool:
				imgui.Checkbox(f.Name, vf.Addr().Interface().(*bool))
			case reflect.Interface:
				q := reflect.ValueOf(vf.Interface())
				name := fmt.Sprintf("%s (%s)", f.Name, q.String())
				if imgui.TreeNodeV(name, imgui.TreeNodeFlagsDefaultOpen) {
					buildStructWidget(q)
					imgui.TreePop()
				}
			default:
				imgui.LabelText(f.Name, fmt.Sprintf("%s isn't not supported", kind.String()))
			}
		}
	default:
		imgui.LabelText("???", fmt.Sprintf("%s isn't not supported", typ.String()))
	}
}

func stringOrDefault(val, def string) string {
	if val != "" {
		return val
	}
	return def
}
