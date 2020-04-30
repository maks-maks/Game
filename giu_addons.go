package main

import "github.com/AllenDang/giu/imgui"

type DragFloatWidget struct {
	label  string
	value  *float32
	speed  float32
	min    float32
	max    float32
	format string
	power  float32
}

func (d *DragFloatWidget) Build() {
	imgui.DragFloatV(d.label, d.value, d.speed, d.min, d.max, d.format, d.power)
}

func DragFloat(label string, value *float32) *DragFloatWidget {
	return DragFloatV(label, value, 1.0, 0, 0, "%.3f", 1)
}

func DragFloatV(label string, value *float32, speed float32, min, max float32, format string, power float32) *DragFloatWidget {
	return &DragFloatWidget{
		label:  label,
		value:  value,
		speed:  speed,
		min:    min,
		max:    max,
		format: format,
		power:  power,
	}
}
