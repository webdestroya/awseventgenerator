package awseventgenerator

import "strings"

type Config struct {
	// A string for the package name
	// or a function that is passed the schema
	PackageName any

	// whether enum types should be actual enums
	// adds a helper .Values() method that lists all possible values
	GenerateEnums bool

	// All values, regardless of requiredness should be pointers
	AlwaysPointerize bool

	// whether to emit marshal/unmarshal functions
	// EmitMarshallers bool

	// whether to skip pretty format the code
	NoFormatCode bool

	// Set the root element name. default is AwsEvent
	RootElement string

	// do not emit the Source/DetailType constants
	SkipEventConstants bool

	// makes very violent marshallers, that will fail if a required field is missing
	EnforceRequiredInMarshallers bool
}

func (c *Config) InitDefaults() {

}

type PackageNameFunc = func(*Schema) string

func DefaultPackageNameFunc(s *Schema) string {
	if s.AwsDetailType != "" {
		return strings.ToLower(s.AwsDetailType)
	}

	if s.Title != "" {
		return strings.ToLower(s.Title)
	}

	return "main"
}
