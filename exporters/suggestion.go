package exporters

import (
	"github.com/moccalotto/nick/field"
	"log"
	"strconv"
)

// Exporter interface
type Exporter interface {
	Export(f *field.Field)
}

type SuggestionExporter struct {
	Vars     map[string]string
	Fallback Exporter
}

func NewSuggestionExporter(v map[string]string, fallback Exporter) *SuggestionExporter {
	return &SuggestionExporter{v, fallback}
}

func (e *SuggestionExporter) image(f *field.Field) {
	ie := NewImageExporter()

	var err error = nil

	if fn, ok := e.Vars["suggestion.export.file"]; ok {
		ie.FileName = fn
	}

	if fmt, ok := e.Vars["suggestion.export.format"]; ok {
		ie.Format = fmt
	}

	if w, ok := e.Vars["suggestion.export.width"]; ok {
		if ie.Width, err = strconv.Atoi(w); err != nil {
			log.Fatalf("Bad value for suggestion export.width: '%s'", w)
		}
	}

	if h, ok := e.Vars["suggestion.export.height"]; ok {
		if ie.Height, err = strconv.Atoi(h); err != nil {
			log.Fatalf("Bad value for suggestion export.height: '%s'", h)
		}
	}

	if scale, ok := e.Vars["suggestion.export.scale"]; ok {
		if ie.Scale, err = strconv.ParseFloat(scale, 64); err != nil {
			log.Fatalf("Bad value for suggestion export.scale: '%s'", scale)
		}
	}

	if al, ok := e.Vars["suggestion.export.algorithm"]; ok {
		ie.Algorithm = al
	}

	ie.Export(f)
}

func (e *SuggestionExporter) Export(f *field.Field) {
	ex, ok := e.Vars["suggestion.export.type"]

	if !ok {
		e.Fallback.Export(f)
		return
	}

	switch ex {
	case "image":
		e.image(f)
	default:
		log.Fatalf("Unknown exporter: %s", ex)
	}
}
