package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

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

func checkerror() {
	if code := gl.GetError(); code != 0 {
		panic(fmt.Sprintf("gl.Error = %d", code))
	}
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
		log.Fatal(err)
	}

	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*4, gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	checkerror()

	program, err := CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Fatal(err)
	}

	gl.ClearColor(0.4, 0.4, 0.4, 1.0)

	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	for !window.ShouldClose() && (window.GetKey(glfw.KeyEscape) != glfw.Press) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.VertexAttribPointer(
			0,        // attribute 0
			3,        // size
			gl.FLOAT, // xtype
			false,    // normalized
			0,        // stride

			gl.PtrOffset(0), // pointer
		)

		gl.DrawArrays(gl.TRIANGLES, 0, 4)
		gl.DisableVertexAttribArray(0)

		checkerror()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	triangleVertices = []float32{
		-1, -1, 0,
		1, -1, 0,
		0, 1, 0,
	}

	vertexShader = `
		#version 330 core

		layout(location = 0) in vec3 vertexPosition_modelSpace;

		void main(){
			gl_Position.xyz = vertexPosition_modelSpace;
			gl_Position.w = 1.0;
		}
	` + "\x00"

	fragmentShader = `
		#version 330 core
		out vec3 color;

		void main(){
			color = vec3(1, 0, 0);
		}
	` + "\x00"
)

func CreateProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))

		return 0, fmt.Errorf("Linking failed: %v", log)
	}

	return program, nil
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))

		return 0, fmt.Errorf("Compiling %v failed: %v", source, log)
	}

	return shader, nil
}
