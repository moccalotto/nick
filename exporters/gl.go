package exporters

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/moccalotto/nick/field"
	"log"
	"strings"
)

// ImageExporter exports images to files
type GlExporter struct {
	// If both Width and Height are 0, the image is sized via Scale instead
	Width         int          // Scale to new width. If 0, aspect ratio is preserved
	Height        int          // Scale to new height.If 0, aspect ratio is preserved.
	Title         string       // title of window
	VertShaderSrc string       // source of vertex shader
	FragShaderSrc string       // source of fragment shader
	CellPoints    []float32    // The points (offsets rather) that make up a cell
	prog          uint32       // handle to the opengl program
	window        *glfw.Window // GLFW window.
	vertArr       uint32       // handle to the vertex array object
	vertCount     int32        // number of vertices
	vertShader    uint32       // handle to compiled vertex shader
	fragShader    uint32       // handle to compiled fragment shader
}

func NewGlExporter(w, h int) *GlExporter {
	return &GlExporter{
		Width:  w,
		Height: h,
		Title:  "Cave",

		// simplest vertex shader
		VertShaderSrc: `
			#version 410
			in vec3 vp;
			void main() {
				gl_Position = vec4(vp, 1.0);
			}` + "\x00",

		// simplest fragment shader. White color
		FragShaderSrc: `
			#version 410
			out vec4 frag_colour;
			void main() {
				frag_colour = vec4(1, 1, 1, 1.0);
			}` + "\x00",
		CellPoints: []float32{
			// Bottom left right-angle triangle
			-1, 1, 0,
			1, -1, 0,
			-1, -1, 0,

			// Top right right-angle triangle
			-1, 1, 0,
			1, 1, 0,
			1, -1, 0,
		},
	}
}

func (e *GlExporter) Export(f *field.Field) {
	defer e.cleanup()

	e.init(f)

	for !e.window.ShouldClose() {
		e.update()
	}
}

func (e *GlExporter) cleanup() {
	if e.window != nil {
		glfw.Terminate()
	}
}

func (e *GlExporter) init(f *field.Field) {
	e.initWindow()
	e.initOpenGL()
	points := e.makePoints(f)
	e.initVertices(points)
}

func (e *GlExporter) makePoints(f *field.Field) []float32 {
	square := []float32{
		// Bottom left right-angle triangle
		-1, 1, 0,
		1, -1, 0,
		-1, -1, 0,

		// Top right right-angle triangle
		-1, 1, 0,
		1, 1, 0,
		1, -1, 0,
	}

	points := []float32{}
	w := f.Width()
	h := f.Height()

	f.Walk(func(x, y int, c field.Cell) {
		if c.Dead() {
			return
		}

		for i := 0; i < len(square); i++ {
			var factor float32
			var size float32
			switch i % 3 {
			case 0:
				// x-coordinates are at 0, 3, 6, etc.
				// stretch cell to fit into viewport,
				// which as width of 1.0
				size = 1.0 / float32(w)
				factor = float32(x) * size
			case 1:
				// y-coordinates are at 1, 4, 7, etc.
				// stretch cell to fit into viewport,
				// which as height of 1.0
				size = 1.0 / float32(h)
				factor = float32(y) * size
			default:
				// z-coordinates are at 2,5,8, etc.
				// and we skip them.
				points = append(points, 0.0)
				continue
			}

			if square[i] < 0 {
				points = append(points, factor*2-1)
			} else {
				points = append(points, ((factor+size)*2)-1)
			}
		}
	})

	return points
}

func (e *GlExporter) initWindow() {
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	wnd, err := glfw.CreateWindow(e.Width, e.Height, "Nick the cave man", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	e.window = wnd
	e.window.MakeContextCurrent()
}

func (e *GlExporter) initOpenGL() {
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	e.vertShader = e.compileShader(&e.VertShaderSrc, gl.VERTEX_SHADER)
	e.fragShader = e.compileShader(&e.FragShaderSrc, gl.FRAGMENT_SHADER)

	e.prog = gl.CreateProgram()
	gl.AttachShader(e.prog, e.vertShader)
	gl.AttachShader(e.prog, e.fragShader)
	gl.LinkProgram(e.prog)
}

func (e *GlExporter) update() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(e.prog)

	gl.BindVertexArray(e.vertArr)
	gl.DrawArrays(gl.TRIANGLES, 0, e.vertCount)

	glfw.PollEvents()
	e.window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func (e *GlExporter) initVertices(points []float32) {
	e.vertCount = int32(len(points)) / 3
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

func (e *GlExporter) compileShader(source *string, shaderType uint32) uint32 {
	shaderHandle := gl.CreateShader(shaderType)

	csources, freeFunc := gl.Strs(*source)
	gl.ShaderSource(shaderHandle, 1, csources, nil)
	freeFunc()
	gl.CompileShader(shaderHandle)

	var status int32
	gl.GetShaderiv(shaderHandle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderHandle, gl.INFO_LOG_LENGTH, &logLength)

		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderHandle, logLength, nil, gl.Str(msg))

		log.Fatalf("failed to compile %v: %v", *source, msg)
	}

	return shaderHandle
}
