package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func drawRect(x, y, w, h int32, color [3]float32) {
	// Triangle
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3i(x, y, 0)
	gl.Vertex3i(x, y+h, 0)
	gl.Vertex3i(x+w, y, 0)
	gl.Vertex3i(x+w, y+h, 0)
	gl.Vertex3i(x, y+h, 0)
	gl.Vertex3i(x+w, y, 0)
	gl.End()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	w, h := 640, 640

	window, err := glfw.CreateWindow(w, h, "2048", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Ortho(0, 640, 640, 0, -1, 1)

	for !window.ShouldClose() {

		gl.ClearColor(1, 1, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		for j := range 4 {
			offset_y := int32(100 * j)
			for i := range 4 {
				offset_x := int32(100 * i)
				drawRect(120+offset_x, 120+offset_y, 90, 90, [3]float32{0.8, 0, 0})
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Millisecond * 10)
	}
}
