package main

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// adopted from https://github.com/g3n/engine/blob/master/texture/Textureanimator.go

type TextureAnimator struct {
	Tex           *texture.Texture2D // pointer to texture being displayed
	DispTime      float32            // disply duration of each tile in milliseconds (default = 1.0/30.0)
	MaxCycles     int                // maximum number of cycles (default = 0 - continuous)
	Cycles        int                // current number of complete cycles
	Columns       int                // number of columns
	Rows          int                // number of rows
	Index         int                // current frame index
	TileDisplayed float32            // time from tile started to be displayed
}

func NewTextureAnimator(tex *texture.Texture2D, htiles, vtiles int) *TextureAnimator {
	a := new(TextureAnimator)
	a.Tex = tex
	a.Columns = htiles
	a.Rows = vtiles
	a.DispTime = 16
	a.MaxCycles = 0

	tex.SetWrapS(gls.REPEAT)
	tex.SetWrapT(gls.REPEAT)
	tex.SetRepeat(1/float32(a.Columns), 1/float32(a.Rows))

	a.Reset()
	return a
}

func (a *TextureAnimator) Clone() interface{} {
	return &TextureAnimator{
		Tex:           a.Tex,
		DispTime:      a.DispTime,
		MaxCycles:     a.MaxCycles,
		Cycles:        a.Cycles,
		Columns:       a.Columns,
		Rows:          a.Rows,
		Index:         a.Index,
		TileDisplayed: a.TileDisplayed,
	}
}

func (a *TextureAnimator) Reset() {
	a.TileDisplayed = 0
	a.Index = 0
	a.Cycles = 0
}

func (a *TextureAnimator) Update(dt float32) {
	if a.MaxCycles > 0 && a.Cycles >= a.MaxCycles {
		return
	}

	a.TileDisplayed += dt
	if a.TileDisplayed >= a.DispTime {
		dFrame := math32.Floor(a.TileDisplayed / a.DispTime)
		a.TileDisplayed -= a.DispTime * dFrame

		// dCycle := math32.Floor(dFrame / (float32(a.Rows * a.Columns)))
		dCycle := (a.Index + int(dFrame)) / (a.Rows * a.Columns)
		a.Cycles += dCycle
		a.Index = (a.Index + int(dFrame)) % (a.Rows * a.Columns)

		if a.MaxCycles > 0 && a.Cycles >= a.MaxCycles {
			a.Index = a.Rows*a.Columns - 1
		}
	}

	iRow := a.Index / a.Columns
	iCol := a.Index % a.Columns

	a.Tex.SetOffset(float32(iCol)/float32(a.Columns), float32(iRow)/float32(a.Rows))
}
