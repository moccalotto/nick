package exporters

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	"github.com/moccalotto/nick/field"
	"github.com/moccalotto/nick/machine"
	"runtime"
)

type ThreeDExporter struct {
	Width   int
	Height  int
	Title   string
	Machine *machine.Machine
}

func NewThreeDExporter(m *machine.Machine) *ThreeDExporter {
	return &ThreeDExporter{
		Width:   1200,
		Height:  800,
		Title:   "Nick",
		Machine: m,
	}
}

func init() {
	runtime.LockOSThread()
}

func (e *ThreeDExporter) Export() error {

	app, _ := application.Create(application.Options{
		Title:  "Nick",
		Width:  1440,
		Height: 900,
	})

	sizeFactor := float32(0.75)
	offX := -float32(e.Machine.Cave.Width()) / 2.0  // width-offset
	offY := float32(0.0)                            // elevation offset
	offZ := -float32(e.Machine.Cave.Height()) / 2.0 // height-offset

	e.Machine.Cave.Walk(func(x, y int, c field.Cell) {
		elevation := float32(c)
		if elevation <= 0.0001 {
			return
		}

		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		geom := geometry.NewBox(sizeFactor, elevation*sizeFactor, sizeFactor)
		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(
			float32(x)+offX,
			sizeFactor*elevation/2+offY,
			float32(y)+offZ,
		)
		app.Scene().Add(mesh)

	})

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 20.0)
	pointLight.SetPosition(2, 10, 2)

	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	// Creates a grid helper and saves its pointer in the test state
	grid := graphic.NewGridHelper(50, 1, &math32.Color{0.4, 0.4, 0.4})
	app.Scene().Add(grid)

	app.CameraPersp().SetPosition(0, 10, 0)
	app.Run()

	return nil
}
