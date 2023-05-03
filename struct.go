package awseventgenerator

import "fmt"

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

	// If this should just be a straight alias (type ThisThing = XXXX)
	AliasType string
}

func (s *Struct) finalize(g *Generator) {

	if len(s.Fields) == 0 {
		if s.AdditionalType == "" {
			s.AliasType = "interface{}"
		} else {
			s.AliasType = fmt.Sprintf("map[string]%s", s.AdditionalType)
		}
	}
}

func (s *Struct) finalizeFields(g *Generator) {
	for k, v := range s.Fields {
		v := v
		v.finalize(g, s)
		s.Fields[k] = v
	}
}
