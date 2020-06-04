// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modified from g3n/engine/app/app-desktop.go
package main

import (
	"fmt"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

type Application struct {
	window.IWindow                    // Embedded GlfwWindow
	keyState       *window.KeyState   // Keep track of keyboard state
	renderer       *renderer.Renderer // Renderer object
	startTime      time.Time          // Application start time
	frameStart     time.Time          // Frame start time
	frameDelta     time.Duration      // Duration of last frame
}

func NewApplication(width, height int, title string) *Application {
	app := new(Application)
	// Initialize window
	err := window.Init(width, height, title)
	if err != nil {
		panic(err)
	}
	app.IWindow = window.Get()
	app.keyState = window.NewKeyState(app)

	app.renderer = renderer.NewRenderer(app.Gls())
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	return app
}

// Run starts the update loop.
// It calls the user-provided update function every frame.
func (a *Application) Run(update func(rend *renderer.Renderer, deltaTime time.Duration)) {
	a.startTime = time.Now()
	a.frameStart = time.Now()

	for true {
		// If Exit() was called or there was an attempt to close the window dispatch OnExit event for subscribers.
		// If no subscriber cancelled the event, terminate the application.
		if a.IWindow.(*window.GlfwWindow).ShouldClose() {
			a.Dispatch(app.OnExit, nil)
			break
		}

		now := time.Now()
		a.frameDelta = now.Sub(a.frameStart)
		a.frameStart = now

		update(a.renderer, a.frameDelta)
		// Swap buffers and poll events
		a.IWindow.(*window.GlfwWindow).SwapBuffers()
		a.IWindow.(*window.GlfwWindow).PollEvents()
	}

	a.Destroy()
}

// Exit requests to terminate the application
// Application will dispatch OnQuit events to registered subscribers which
// can cancel the process by calling CancelDispatch().
func (a *Application) Exit() {
	a.IWindow.(*window.GlfwWindow).SetShouldClose(true)
}

// Renderer returns the application's renderer.
func (a *Application) Renderer() *renderer.Renderer {
	return a.renderer
}

// KeyState returns the application's KeyState.
func (a *Application) KeyState() *window.KeyState {
	return a.keyState
}

// RunTime returns the elapsed duration since the call to Run().
func (a *Application) RunTime() time.Duration {
	return time.Now().Sub(a.startTime)
}
