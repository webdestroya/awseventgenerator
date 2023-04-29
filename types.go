package awseventgenerator

import "strings"

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
}

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
}

type Enum struct {
	Name   string
	Values []string
}

func (e Enum) id() string {
	return e.Name + ":" + strings.Join(e.Values, "/")
}

const (
	eventSourceConstName     = "Source"
	eventDetailTypeConstName = "DetailType"
)
