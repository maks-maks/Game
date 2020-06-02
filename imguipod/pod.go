package imguipod

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go/v2"
)

func (pod *ImguiPod) BeginImgui() {
	pod.Platform.NewFrame()
	imgui.NewFrame()
}

func (pod *ImguiPod) EndImgui() {
	imgui.Render()
	pod.Renderer.Render(pod.Platform.DisplaySize(), pod.Platform.DisplaySize(), imgui.RenderedDrawData())
	pod.Platform.PostRender()
}

type ImguiPod struct {
	IO       imgui.IO
	Platform Platform
	Renderer *OpenGL3
}

func New(w *glfw.Window) *ImguiPod {
	pod := &ImguiPod{}

	err := imgui.CreateContext(nil).SetCurrent()
	if err != nil {
		panic(err)
	}

	io := imgui.CurrentIO()

	pod.IO = io

	pod.Renderer, err = NewOpenGL3(io)
	if err != nil {
		panic(err)
	}

	pod.Platform, err = NewGLFWFromWindow(io, w)
	if err != nil {
		panic(err)
	}

	return pod
}
