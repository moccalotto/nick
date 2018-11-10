package exporters

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
)

type ThreeDExporter struct {
	Width  int
	Height int
	Title  string
}

func NewThreeDExporter() *ThreeDExporter {
	return &ThreeDExporter{
		Width:  1440,
		Height: 900,
		Title:  "Nick",
	}
}

func (e *ThreeDExporter) Export() error {

	app, _ := application.Create(application.Options{
		Title:  "Nick",
		Width:  1440,
		Height: 900,
	})

	// Create a blue torus and add it to the scene
	geom := geometry.NewBox(1.0, 1.0, 1.0)
	mat := material.NewPhong(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	app.Scene().Add(mesh)

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)

	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	app.CameraPersp().SetPosition(0, 0, 3)
	return app.Run()
}
