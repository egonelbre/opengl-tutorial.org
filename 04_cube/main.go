package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"

	_ "image/png"

	"github.com/kardianos/osext"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

	program, err := CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Fatal(err)
	}
	gl.UseProgram(program)

	Projection := mgl32.Perspective(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	View := mgl32.LookAt(
		4, 3, 3,
		0, 0, 0,
		0, 1, 0,
	)
	Model := mgl32.Ident4()

	MVP := Projection.Mul4(View).Mul4(Model)

	MVP_ID := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	gl.UniformMatrix4fv(MVP_ID, 1, false, &MVP[0])

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	var cbo uint32
	gl.GenBuffers(1, &cbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeColors)*4, gl.Ptr(cubeColors), gl.STATIC_DRAW)

	checkerror()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.4, 0.4, 0.4, 1.0)

	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	for !window.ShouldClose() && (window.GetKey(glfw.KeyEscape) != glfw.Press) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		for i := range cubeVertices {
			cubeVertices[i] += (rand.Float32() - 0.5) / 100
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, cbo)
		for i := range cubeColors {
			cubeColors[i] += (rand.Float32() - 0.5) / 10
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(cubeColors)*4, gl.Ptr(cubeColors), gl.STATIC_DRAW)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.DrawArrays(gl.TRIANGLES, 0, 12*3)
		gl.DisableVertexAttribArray(0)

		checkerror()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	cubeVertices = []float32{
		-1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0,

		1.0, -1.0, 1.0, -1.0, -1.0, -1.0, 1.0, -1.0, -1.0,
		1.0, 1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, -1.0,

		-1.0, -1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0,
		1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, -1.0,

		-1.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0,
		1.0, 1.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0,

		1.0, -1.0, -1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0,

		1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0,
		1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, -1.0, 1.0,
	}

	cubeColors = []float32{
		0.583, 0.771, 0.014, 0.609, 0.115, 0.436, 0.327, 0.483, 0.844,
		0.822, 0.569, 0.201, 0.435, 0.602, 0.223, 0.310, 0.747, 0.185,

		0.597, 0.770, 0.761, 0.559, 0.436, 0.730, 0.359, 0.583, 0.152,
		0.483, 0.596, 0.789, 0.559, 0.861, 0.639, 0.195, 0.548, 0.859,

		0.014, 0.184, 0.576, 0.771, 0.328, 0.970, 0.406, 0.615, 0.116,
		0.676, 0.977, 0.133, 0.971, 0.572, 0.833, 0.140, 0.616, 0.489,

		0.997, 0.513, 0.064, 0.945, 0.719, 0.592, 0.543, 0.021, 0.978,
		0.279, 0.317, 0.505, 0.167, 0.620, 0.077, 0.347, 0.857, 0.137,

		0.055, 0.953, 0.042, 0.714, 0.505, 0.345, 0.783, 0.290, 0.734,
		0.722, 0.645, 0.174, 0.302, 0.455, 0.848, 0.225, 0.587, 0.040,

		0.517, 0.713, 0.338, 0.053, 0.959, 0.120, 0.393, 0.621, 0.362,
		0.673, 0.211, 0.457, 0.820, 0.883, 0.371, 0.982, 0.099, 0.879,
	}

	vertexShader = `
		#version 330 core
		layout(location = 0) in vec3 vertexPosition_modelSpace;
		layout(location = 1) in vec3 vertexColor;
		uniform mat4 MVP;

		out vec3 fragmentColor;

		void main(){
			vec4 v = vec4(vertexPosition_modelSpace, 1);
			gl_Position = MVP * v;
			fragmentColor = vertexColor;
		}
	` + "\x00"

	fragmentShader = `
		#version 330 core
		in vec3 fragmentColor;
		out vec3 color;

		void main(){
			color = fragmentColor;
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
