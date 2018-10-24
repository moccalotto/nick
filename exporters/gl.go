package exporters

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/moccalotto/nick/field"
	"log"
	"runtime"
)

// ImageExporter exports images to files
type GlExporter struct {
	// If both Width and Height are 0, the image is sized via Scale instead
	Width     int          // Scale to new width. If 0, aspect ratio is preserved
	Height    int          // Scale to new height.If 0, aspect ratio is preserved.
	Title     string       // title of window
	prog      uint32       // handle to the opengl program
	window    *glfw.Window // GLFW window.
	vertArr   uint32       // handle to the vertex array object
	vertCount uint32       // number of vertices
}

func NewGlExporter(w, h int) *GlExporter {
	return &GlExporter{
		Width:  w,
		Height: h,
		Title:  "Cave",
	}
}

func (e *GlExporter) Export(f *field.Field) {
	runtime.LockOSThread()
	e.init(f)
	e.initWindow()
	defer glfw.Terminate()
	e.initOpenGL()

	for !e.window.ShouldClose() {
		e.update()
	}
}

func (e *GlExporter) init(f *field.Field) {
	e.initWindow()
	e.initOpenGL()
	triangles := e.makePoints(f)
	e.initVertices(points)
}

func (e *GlExporter) makePoints(f *field.Field) {

	// TODO: do we even know the length of this array now?
	// if we only draw living cells, that number is unknown
	// until we've actually counted the living cells.
	points := make([]float32, f.Height()*f.Width()*3)

	halfX := float32(f.Width()) / 2.0
	halfY := float32(f.Hieght()) / 2.0

	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			// DO we only draw living cells?
			// do we draw 2d or 3d ?
			// dafuq?
		}
	}
}

func (e *GlExporter) initWindow() {
	if err := glfw.Init(); err != nil {
		log.Fatalf(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	e.window, err = glfw.CreateWindow(e.Width, e.Height, "Conway's Game of Life", nil, nil)
	if err != nil {
		log.Fatalf(err)
	}
	e.window.MakeContextCurrent()
}

func (e *GlExporter) initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	e.programHandle = gl.CreateProgram()
	gl.LinkProgram(e.prog)
}

func (e *GlExporter) update(f *field.Field) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(e.prog)

	gl.BindVertexArray(e.vertArr)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func (e *GlExporter) initVertices(points []float32) {
	e.vertCount = len(points)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &e.vertArr)
	gl.BindVertexArray(e.vertArr)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}
