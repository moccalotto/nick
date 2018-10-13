package exporters

import (
	"github.com/moccalotto/nick/field"
)

// Exporter interface
type Exporter interface {
	Export(f *field.Field)
}
