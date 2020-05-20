package main

import (
	"fmt"
	"reflect"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/bytearena/ecs"
)

var demo = false
var curEntityID ecs.EntityID = 0
var delayMs int32 = 30

func leftPanel() *g.Layout {
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

	return &g.Layout{
		g.ListBox("Entities", items,
			func(i int) {
				curEntityID = q.Entities()[i].ID
			},
			func(selectedIndex int) {}),
	}
}

func rightPanel() *g.Layout {
	q := ecsManager.GetEntityByID(curEntityID)
	if q == nil {
		return &g.Layout{}
	}

	l := make(g.Layout, 0)

	for _, component := range ecsManager.components {
		if q.Entity.HasComponent(component) {
			n := g.TreeNode(ecsManager.ComponentName(component), g.TreeNodeFlagsDefaultOpen, entityComponentLayout(q.Entity, component))
			l = append(l, n)
		}
	}

	return &l
}

func entityComponentLayout(e *ecs.Entity, c *ecs.Component) g.Layout {
	d, _ := e.GetComponentData(c)

	val := reflect.ValueOf(d)
	val = reflect.Indirect(val)

	return structLayout(val)
}

func structLayout(val reflect.Value) g.Layout {
	val = reflect.Indirect(val)
	typ := val.Type()

	l := make(g.Layout, 0)

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
				w := g.InputText(f.Name, 0, vf.Addr().Interface().(*string))
				l = append(l, w)
			case reflect.Float32:
				format = stringOrDefault(format, "%.3f")
				w := DragFloatV(f.Name, vf.Addr().Interface().(*float32), 1, 0, 0, format, 1)
				l = append(l, w)
			case reflect.Int32:
				format = stringOrDefault(format, "%d")
				w := g.DragIntV(f.Name, vf.Addr().Interface().(*int32), 1.0, 0, 0, format)
				l = append(l, w)
			case reflect.Uint32:
				format = stringOrDefault(format, "%d")
				// imgui.Drag
				// w := g.DragIntV(f.Name, vf.Addr().Interface().(*int32), 1.0, 0, 0, format)
				w := LabelText(f.Name, fmt.Sprintf("%d", vf.Interface().(ecs.EntityID)))
				l = append(l, w)
			case reflect.Bool:
				w := g.Checkbox(f.Name, vf.Addr().Interface().(*bool), nil)
				l = append(l, w)
			case reflect.Interface:
				q := reflect.ValueOf(vf.Interface())
				name := fmt.Sprintf("%s (%s)", f.Name, q.String())
				n := g.TreeNode(name, g.TreeNodeFlagsDefaultOpen, structLayout(q))
				l = append(l, n)
			default:
				w := LabelText(f.Name, fmt.Sprintf("%s isn't not supported", kind.String()))
				l = append(l, w)
			}
		}
	default:
	}

	return l
}

func stringOrDefault(val, def string) string {
	if val != "" {
		return val
	}
	return def
}
