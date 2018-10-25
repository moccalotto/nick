package exporters

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/moccalotto/nick/field"
	"log"
)

/**
Note to self:
OpenGL uses a right-handed.coordinate system.
This means that:
	when X increases, you move to the left,
	when Y increases, you move up (which is inverse of how images are rendered, but matches how math graphs are normally displayed),
	when Z increases, you move "towards the camera". This might be the opposite of what one would expect of a Z-axis.

I want to map the "depth" of my field (i.e. the tallnes or value of the cell) to the Z-axis, so no issue there
I want to map the width of my field to the X-axis, so no issue there either
I want to map the height of my field to the Y-axis. I need to flip the Y coordinate so that Y_new := Height - Y_original

Use this example: https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
Also: http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
*/

// ImageExporter exports images to files
type GlExporter struct {
	// If both Width and Height are 0, the image is sized via Scale instead
	Width         int          // Scale to new width. If 0, aspect ratio is preserved
	Height        int          // Scale to new height.If 0, aspect ratio is preserved.
	Title         string       // title of window
	VertShaderSrc string       // source of vertex shader
	FragShaderSrc string       // source of fragment shader
	prog          uint32       // handle to the opengl program
	window        *glfw.Window // GLFW window.
	vertArr       uint32       // handle to the vertex array object
	vertCount     int32        // number of vertices
	vertShader    uint32       // handle to compiled vertex shader
	fragShader    uint32       // handle to compiled fragment shader
	prevTickAt    float64      // time of last tick (glfw value)
}

func NewGlExporter(w, h int) *GlExporter {
	return &GlExporter{
		Width:  w,
		Height: h,
		Title:  "Cave",

		// simplest vertex shader
		VertShaderSrc: `
			#version 410

			uniform mat4 projection;
			uniform mat4 camera;
			uniform mat4 model;
			in vec3 vert;
			in vec2 vertTexCoord;
			out vec2 fragTexCoord;
			void main() {
				fragTexCoord = vertTexCoord;
				gl_Position = projection * camera * model * vec4(vert, 1);
			}
			` + "\x00",

		// simplest fragment shader. White color
		FragShaderSrc: `
			#version 410
			uniform sampler2D tex;
			in vec2 fragTexCoord;
			out vec4 outputColor;
			void main() {
			    outputColor = texture(tex, fragTexCoord);
			}
			` + "\x00",
	}
}

func (e *GlExporter) Export(f *field.Field) {
	defer e.cleanup()

	e.init(f)

	e.prevTickAt = glfw.GetTime()

	for !e.window.ShouldClose() {
		e.tick()
	}
}

func (e *GlExporter) cleanup() {
	if e.window != nil {
		glfw.Terminate()
	}
}

func (e *GlExporter) init(f *field.Field) {
	e.initGlfw()   // create a window
	e.initOpenGL() // start opengl
}

func (e *GlExporter) initGlfw() {
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
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

	e.prog = gl.CreateProgram()
	e.vertShader = e.loadShader(&e.VertShaderSrc, gl.VERTEX_SHADER)
	e.fragShader = e.loadShader(&e.FragShaderSrc, gl.FRAGMENT_SHADER)
	gl.LinkProgram(e.prog)
}

func (e *GlExporter) tick() {
	// clear
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// updates
	now := glfw.GetTime()
	elapsed := now - e.prevTickAt
	e.prevTickAt = now
	angle += elapsed
	/* model = */
	mgl.HomogRotate3D(float32(angle), mgl.Vec3{0, 1, 0})

	// render
	gl.UseProgram(e.prog)

	e.window.SwapBuffers()
	glfw.PollEvents()
}

func (e *GlExporter) loadShader(source *string, shaderType uint32) {
	if e.prog == 0 {
		log.Fatal("Program not initialized yet! Besure to call CreateProgram")
	}

	shaderHandle := gl.CreateShader(shaderType)

	src, free := gl.Strs(*source)
	defer free()

	gl.ShaderSource(shaderHandle, 1, src, nil)
	gl.CompileShader(shaderHandle)

	var success int32
	gl.GetShaderiv(shaderHandle, gl.COMPILE_STATUS, &success)

	if success == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shaderHandle, gl.INFO_LOG_LENGTH, &logLen)

		infoLog := make([]byte, logLen)
		gl.GetShaderInfoLog(shaderHandle, logLen, nil, &infoLog[0])
		log.Fatalf("error compiling vertex shader: %s", string(infoLog))
	}

	gl.AttachShader(e.prog, shaderHandle)
}
