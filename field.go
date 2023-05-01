package awseventgenerator

// Field defines the data required to generate a field in Go.
type Field struct {
	// The golang name, e.g. "Address1"
	Name string
	// The JSON name, e.g. "address1"
	JSONName string
	// The golang type of the field, e.g. a built-in type like "string" or the name of a struct generated
	// from the JSON schema.
	Type string
	// Required is set to true when the field is required.
	Required    bool
	Description string
	EnumValues  []string
	Format      string

	// the actual go type to write to the file
	FinalType string

	finalized bool
}

func (f *Field) finalize(g *Generator, s *Struct) {
	f.finalized = true
}
