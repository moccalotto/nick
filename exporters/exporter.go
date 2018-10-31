package exporters

// Exporter interface
type Exporter interface {
	Export() error
}
