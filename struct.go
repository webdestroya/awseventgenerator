package awseventgenerator

// Struct defines the data required to generate a struct in Go.
type Struct struct {
	// The ID within the JSON schema, e.g. #/definitions/address
	ID string
	// The golang name, e.g. "Address"
	Name string
	// Description of the struct
	Description string
	Fields      map[string]Field

	GenerateCode   bool
	AdditionalType string
	importTypes    []string

	EnumValues []string
	IsEnum     bool
}

func (s *Struct) finalize(g *Generator) {
	for k, v := range s.Fields {
		v := v
		v.finalize(g, s)
		s.Fields[k] = v
	}
}
