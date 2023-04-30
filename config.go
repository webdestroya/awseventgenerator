package awseventgenerator

import "strings"

type Config struct {
	// A string for the package name
	// or a function that is passed the schema
	PackageName any

	// whether enum types should be actual enums
	GenerateEnums bool

	// adds a helper .Values() method that lists all possible values
	GenerateEnumValueMethod bool

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

	// Adds helpers for running tests, not meant to be included normally
	AddTestHelpers bool
}

func (c *Config) InitDefaults() {
	if c.RootElement == "" {
		c.RootElement = "AwsEvent"
	}

	if c.PackageName == nil {
		c.PackageName = DefaultPackageNameFunc
	}

	if c.GenerateEnumValueMethod {
		c.GenerateEnums = true
	}
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
