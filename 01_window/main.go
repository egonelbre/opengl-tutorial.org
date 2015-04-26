package main

import (
	"fmt"
	"log"
	"runtime"

	_ "image/png"

	"github.com/kardianos/osext"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const WindowWidth = 800
const WindowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Initialize Window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	exe, _ := osext.Executable()
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, exe, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)

	for !window.ShouldClose() && (window.GetKey(glfw.KeyEscape) != glfw.Press) {
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
