package field

type Modifier interface {
	ApplyToField(f *Field)
}

func (f *Field) Apply(m Modifier) {
	m.ApplyToField(f)
}
