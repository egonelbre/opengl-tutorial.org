package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/kardianos/osext"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/egonelbre/opengl-tutorial.org/dds"
	"github.com/egonelbre/opengl-tutorial.org/obj"
	"github.com/egonelbre/opengl-tutorial.org/shaders"
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

type Controls struct {
	Window *glfw.Window

	Projection mgl32.Mat4
	Camera     mgl32.Mat4

	Position        mgl32.Vec3
	HorizontalAngle float32
	VerticalAngle   float32
	InitialFoV      float32

	Speed      float32
	MouseSpeed float32
}

func NewControls(window *glfw.Window) *Controls {
	return &Controls{
		Window: window,

		Projection: mgl32.Perspective(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0),
		Camera: mgl32.LookAt(
			4, 3, 3,
			0, 0, 0,
			0, 1, 0,
		),

		Position:        mgl32.Vec3{0, 0, 5},
		HorizontalAngle: math.Pi,
		VerticalAngle:   0.0,
		InitialFoV:      45.0,

		Speed:      3.0,
		MouseSpeed: 0.05,
	}
}

func Cos(v float32) float32 { return float32(math.Cos(float64(v))) }
func Sin(v float32) float32 { return float32(math.Sin(float64(v))) }

func (c *Controls) Update(dt float32) {
	W := c.Window
	mouseX, mouseY := W.GetCursorPos()
	W.SetCursorPos(float64(WindowWidth/2), float64(WindowHeight/2))

	c.HorizontalAngle += c.MouseSpeed * dt * float32(WindowWidth/2-mouseX)
	c.VerticalAngle += c.MouseSpeed * dt * float32(WindowHeight/2-mouseY)

	direction := mgl32.Vec3{
		Cos(c.VerticalAngle) * Sin(c.HorizontalAngle),
		Sin(c.VerticalAngle),
		Cos(c.VerticalAngle) * Cos(c.HorizontalAngle),
	}

	right := mgl32.Vec3{
		Sin(c.HorizontalAngle - math.Pi/2),
		0,
		Cos(c.HorizontalAngle - math.Pi/2),
	}
	up := right.Cross(direction)

	if W.GetKey(glfw.KeyUp) == glfw.Press || W.GetKey(glfw.KeyW) == glfw.Press {
		c.Position = c.Position.Add(direction.Mul(dt).Mul(c.Speed))
	}
	if W.GetKey(glfw.KeyDown) == glfw.Press || W.GetKey(glfw.KeyS) == glfw.Press {
		c.Position = c.Position.Sub(direction.Mul(dt).Mul(c.Speed))
	}
	if W.GetKey(glfw.KeyRight) == glfw.Press || W.GetKey(glfw.KeyD) == glfw.Press {
		c.Position = c.Position.Add(right.Mul(dt).Mul(c.Speed))
	}
	if W.GetKey(glfw.KeyLeft) == glfw.Press || W.GetKey(glfw.KeyA) == glfw.Press {
		c.Position = c.Position.Sub(right.Mul(dt).Mul(c.Speed))
	}

	c.Projection = mgl32.Perspective(c.InitialFoV, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	c.Camera = mgl32.LookAtV(c.Position, c.Position.Add(direction), up)
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

	program, err := shaders.Load("transform.vert", "texture.frag")
	if err != nil {
		log.Fatal(err)
	}
	gl.UseProgram(program)

	ProjectionID := gl.GetUniformLocation(program, gl.Str("Projection\x00"))
	CameraID := gl.GetUniformLocation(program, gl.Str("Camera\x00"))
	ModelID := gl.GetUniformLocation(program, gl.Str("Model\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	texture, err := dds.LoadFile("cube.dds")
	if err != nil {
		log.Fatal(err)
	}

	checkerror()

	// Load Model
	data, err := obj.LoadFile("cube.obj")
	if err != nil {
		log.Fatal(err)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data.Vertex)*4, gl.Ptr(data.Vertex), gl.STATIC_DRAW)

	vertex := uint32(gl.GetAttribLocation(program, gl.Str("vertex\x00")))
	gl.EnableVertexAttribArray(vertex)
	gl.VertexAttribPointer(vertex, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	checkerror()

	var uvbo uint32
	gl.GenBuffers(1, &uvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data.UV)*4, gl.Ptr(data.UV), gl.STATIC_DRAW)

	vertexUV := uint32(gl.GetAttribLocation(program, gl.Str("UV\x00")))
	gl.EnableVertexAttribArray(vertexUV)
	gl.VertexAttribPointer(vertexUV, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	checkerror()

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.4, 0.4, 0.4, 1.0)

	lastTime := glfw.GetTime()

	Model := mgl32.Ident4()
	controls := NewControls(window)
	angle := 0.0

	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	for !window.ShouldClose() && (window.GetKey(glfw.KeyEscape) != glfw.Press) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		deltaTime := time - lastTime
		lastTime = time

		controls.Update(float32(deltaTime))

		Model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		angle += 0.01

		gl.UseProgram(program)

		gl.UniformMatrix4fv(ProjectionID, 1, false, &controls.Projection[0])
		gl.UniformMatrix4fv(CameraID, 1, false, &controls.Camera[0])
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
)
