package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

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
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	title, _ := osext.Executable()
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, title, nil, nil)
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
	Camera := mgl32.LookAt(
		4, 3, 3,
		0, 0, 0,
		0, 1, 0,
	)
	Model := mgl32.Ident4()

	ProjectionID := gl.GetUniformLocation(program, gl.Str("Projection\x00"))
	CameraID := gl.GetUniformLocation(program, gl.Str("Camera\x00"))
	ModelID := gl.GetUniformLocation(program, gl.Str("Model\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeData)*4, gl.Ptr(cubeData), gl.STATIC_DRAW)

	vertex := uint32(gl.GetAttribLocation(program, gl.Str("vertex\x00")))
	gl.EnableVertexAttribArray(vertex)
	gl.VertexAttribPointer(vertex, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	vertexUV := uint32(gl.GetAttribLocation(program, gl.Str("vertexUV\x00")))
	gl.EnableVertexAttribArray(vertexUV)
	gl.VertexAttribPointer(vertexUV, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	texture, err := CreateTexture("cube.png")
	if err != nil {
		log.Fatal(err)
	}

	checkerror()

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.4, 0.4, 0.4, 1.0)

	lastTime := glfw.GetTime()

	Position := mgl32.Vec3{0, 0, 5}
	HorizontalAngle := math.Pi
	VerticalAngle := 0.0
	InitialFoV := float32(45.0)

	Speed := 3.0
	MouseSpeed := 0.05

	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	for !window.ShouldClose() && (window.GetKey(glfw.KeyEscape) != glfw.Press) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		deltaTime := time - lastTime
		deltaTime32 := float32(deltaTime)
		lastTime = time
		speed32 := float32(Speed)

		mouseX, mouseY := window.GetCursorPos()
		window.SetCursorPos(float64(WindowWidth/2), float64(WindowHeight/2))

		HorizontalAngle += MouseSpeed * deltaTime * float64(WindowWidth/2-mouseX)
		VerticalAngle += MouseSpeed * deltaTime * float64(WindowHeight/2-mouseY)

		direction := mgl32.Vec3{
			float32(math.Cos(VerticalAngle) * math.Sin(HorizontalAngle)),
			float32(math.Sin(VerticalAngle)),
			float32(math.Cos(VerticalAngle) * math.Cos(HorizontalAngle)),
		}

		right := mgl32.Vec3{
			float32(math.Sin(HorizontalAngle - math.Pi/2)),
			float32(0),
			float32(math.Cos(HorizontalAngle - math.Pi/2)),
		}
		up := right.Cross(direction)

		if window.GetKey(glfw.KeyUp) == glfw.Press {
			Position = Position.Add(direction.Mul(deltaTime32).Mul(speed32))
		}
		if window.GetKey(glfw.KeyDown) == glfw.Press {
			Position = Position.Sub(direction.Mul(deltaTime32).Mul(speed32))
		}
		if window.GetKey(glfw.KeyRight) == glfw.Press {
			Position = Position.Add(right.Mul(deltaTime32).Mul(speed32))
		}
		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			Position = Position.Sub(right.Mul(deltaTime32).Mul(speed32))
		}

		Projection = mgl32.Perspective(InitialFoV, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
		Camera = mgl32.LookAtV(Position, Position.Add(direction), up)
		// Model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		gl.UseProgram(program)

		gl.UniformMatrix4fv(ProjectionID, 1, false, &Projection[0])
		gl.UniformMatrix4fv(CameraID, 1, false, &Camera[0])
		gl.UniformMatrix4fv(ModelID, 1, false, &Model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 12*3)

		checkerror()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	cubeData = []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,

		// Top
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,

		// Left
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,

		// Right
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}

	vertexShader = `
		#version 330 core
		uniform mat4 Projection;
		uniform mat4 Camera;
		uniform mat4 Model;

		in vec3 vertex;
		in vec2 vertexUV;

		out vec2 UV;

		void main(){
			gl_Position = Projection * Camera * Model * vec4(vertex, 1);
			UV = vertexUV;
		}
	` + "\x00"

	fragmentShader = `
		#version 330 core
		in vec2 UV;
		out vec3 color;
		uniform sampler2D sampler;

		void main(){
			color = texture(sampler, UV).rgb;
		}
	` + "\x00"
)
