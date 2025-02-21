package awseventgenerator

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

// Generator will produce structs from the JSON schema.
type Generator struct {
	schemas  []*Schema
	resolver *RefResolver
	Structs  map[string]Struct
	Aliases  map[string]Field
	// cache for reference types; k=url v=type
	refs      map[string]string
	anonCount int
	config    *Config
	Constants map[string]any

	finalFieldTypes map[string]string

	fileSet *token.FileSet
}

// New creates an instance of a generator which will produce structs.
func NewMulti(config *Config, schemas ...*Schema) *Generator {
	config.InitDefaults()
	return &Generator{
		schemas:         schemas,
		config:          config,
		resolver:        NewRefResolver(schemas),
		Structs:         make(map[string]Struct),
		Aliases:         make(map[string]Field),
		refs:            make(map[string]string),
		Constants:       make(map[string]any),
		finalFieldTypes: make(map[string]string),
		fileSet:         token.NewFileSet(),
	}
}

func New(config *Config, schema *Schema) *Generator {
	schemas := []*Schema{schema}
	return NewMulti(config, schemas...)
}

func (g *Generator) finalize() {

	// First pass, finalize structs
	for k, v := range g.Structs {
		v := v
		v.finalize(g)
		g.Structs[k] = v
	}

	// now go back and do the fields
	for k, v := range g.Structs {
		v := v
		v.finalizeFields(g)
		g.Structs[k] = v

		for _, fv := range v.Fields {
			g.finalFieldTypes[v.Name+":"+fv.Name] = g.determineFieldFinalType(fv)
		}
	}
}

func (g *Generator) determineFieldFinalType(f Field) string {

	if f.Format == "raw" {
		return "json.RawMessage"
	}

	ftype := f.Type
	if ftype == "int" {
		// ftype = "int64"
		ftype = "float64"
	} else if ftype == "string" && f.Format == "date-time" {
		ftype = "time.Time"
	}

	wantsPointer := !f.Required || g.config.AlwaysPointerize

	nonPtrType := strings.TrimPrefix(ftype, "*")

	if _, ok := g.Aliases[nonPtrType]; ok {
		return nonPtrType
	}

	if s, ok := g.Structs[nonPtrType]; ok {
		if s.IsEnum || s.AliasType != "" {
			return nonPtrType
		}
	}

	if wantsPointer && !strings.HasPrefix(ftype, "*") && isPointerable(g, ftype) {
		ftype = "*" + ftype
	}

	if !isPointerable(g, nonPtrType) && strings.HasPrefix(ftype, "*") {
		ftype = nonPtrType
	}

	return ftype
}

func isPointerable(g *Generator, typ string) bool {

	if typ == "interface{}" || strings.HasPrefix(typ, "*") || strings.HasPrefix(typ, "[]") || strings.HasPrefix(typ, "map") {
		return false
	}

	nonPtr := strings.TrimPrefix(typ, "*")

	if a, ok := g.Aliases[nonPtr]; ok {
		if strings.HasPrefix(a.Type, "map") {
			return false
		}
	}

	if s, ok := g.Structs[nonPtr]; ok {
		if s.IsEnum || s.AliasType != "" {
			return false
		}
	}

	return true
}

// Generate creates types from the JSON schemas, keyed by the golang name.
func (g *Generator) Generate() (err error) {
	if err := g.resolver.Init(); err != nil {
		return err
	}

	// extract the types
	for _, schema := range g.schemas {
		name := g.getSchemaName("", schema)
		rootType, err := g.processSchema(name, schema)
		if err != nil {
			return err
		}

		if schema.AwsSource != "" {
			g.Constants[eventSourceConstName] = schema.AwsSource
		}
		if schema.AwsDetailType != "" {
			g.Constants[eventDetailTypeConstName] = schema.AwsDetailType
		}

		// ugh: if it was anything but a struct the type will not be the name...
		if rootType != "*"+name {
			a := Field{
				Name:        name,
				JSONName:    "",
				Type:        rootType,
				Required:    false,
				Description: schema.Description,
				Format:      schema.Format,
				EnumValues:  schema.EnumValues,
			}
			g.Aliases[a.Name] = a
		}
	}

	g.finalize()

	return
}

// process a block of definitions
func (g *Generator) processDefinitions(schema *Schema) error {
	for key, subSchema := range schema.Definitions {
		if _, err := g.processSchema(getGolangName(key), subSchema); err != nil {
			return err
		}
	}
	return nil
}

// process a reference string
func (g *Generator) processReference(schema *Schema) (string, error) {
	schemaPath := g.resolver.GetPath(schema)
	if schema.Reference == "" {
		return "", errors.New("processReference empty reference: " + schemaPath)
	}
	refSchema, err := g.resolver.GetSchemaByReference(schema)
	if err != nil {
		return "", errors.New("processReference: reference \"" + schema.Reference + "\" not found at \"" + schemaPath + "\"")
	}
	if refSchema.GeneratedType == "" {
		// reference is not resolved yet. Do that now.
		refSchemaName := g.getSchemaName("", refSchema)
		typeName, err := g.processSchema(refSchemaName, refSchema)
		if err != nil {
			return "", err
		}
		return typeName, nil
	}
	return refSchema.GeneratedType, nil
}

// returns the type referred to by schema after resolving all dependencies
func (g *Generator) processSchema(schemaName string, schema *Schema) (typ string, err error) {
	if len(schema.Definitions) > 0 {
		if err := g.processDefinitions(schema); err != nil {
			return "", err
		}
	}
	schema.FixMissingTypeValue()
	// if we have multiple schema types, the golang type will be interface{}
	typ = unknownVariableType
	types, isMultiType := schema.MultiType(g.config)
	if len(types) > 0 {
		for _, schemaType := range types {
			name := schemaName
			if isMultiType {
				name = name + "_" + schemaType
			}
			switch schemaType {
			case "object":
				rv, err := g.processObject(name, schema)
				if err != nil {
					return "", err
				}
				if !isMultiType {
					return rv, nil
				}
			case "array":
				rv, err := g.processArray(name, schema)
				if err != nil {
					return "", err
				}
				if !isMultiType {
					return rv, nil
				}
			case "enum":
				rv, err := g.processEnum(name, schema)
				if err != nil {
					return "", err
				}
				if !isMultiType {
					return rv, nil
				}
			default:
				rv, err := getPrimitiveTypeName(schemaType, "", false)
				if err != nil {
					return "", err
				}
				if !isMultiType {
					return rv, nil
				}
			}
		}
	} else {
		if schema.Reference != "" {
			return g.processReference(schema)
		}
	}
	return // return interface{}
}

// name: name of this array, usually the js key
// schema: items element
func (g *Generator) processArray(name string, schema *Schema) (typeStr string, err error) {
	if schema.Items != nil {
		// subType: fallback name in case this array contains inline object without a title
		subName := g.getSchemaName(name+"Items", schema.Items)
		subTyp, err := g.processSchema(subName, schema.Items)
		if err != nil {
			return "", err
		}
		finalType, err := getPrimitiveTypeName("array", subTyp, true)
		if err != nil {
			return "", err
		}
		// only alias root arrays
		if schema.Parent == nil {
			array := Field{
				Name:        name,
				JSONName:    "",
				Type:        finalType,
				Required:    contains(schema.Required, name),
				Description: schema.Description,
				Format:      schema.Format,
				EnumValues:  schema.EnumValues,
			}
			g.Aliases[array.Name] = array
		}
		return finalType, nil
	}
	return "[]" + unknownVariableType, nil
}

func (g *Generator) processEnum(name string, schema *Schema) (typ string, err error) {

	enumName := g.getSchemaName(name, schema) + "Type"

	if s, ok := g.Structs[enumName]; ok {
		if !slices.Equal(s.EnumValues, schema.EnumValues) {
			s.EnumValues = uniqueNonEmptyElementsOf(append(s.EnumValues, schema.EnumValues...))
			g.Structs[enumName] = s
		}
		return getPrimitiveTypeName("object", enumName, false)
	}

	strct := Struct{
		ID:          schema.ID(),
		Name:        enumName,
		Description: schema.Description,
		IsEnum:      true,
		EnumValues:  schema.EnumValues,
		Fields:      make(map[string]Field, 0),
	}

	g.Structs[strct.Name] = strct

	return getPrimitiveTypeName("object", enumName, false)
}

// name: name of the struct (calculated by caller)
// schema: detail incl properties & child objects
// returns: generated type
func (g *Generator) processObject(name string, schema *Schema) (typ string, err error) {
	strct := Struct{
		ID:          schema.ID(),
		Name:        name,
		Description: schema.Description,
		Fields:      make(map[string]Field, len(schema.Properties)),
	}
	// cache the object name in case any sub-schemas recursively reference it
	schema.GeneratedType = "*" + name
	// regular properties
	for propKey, prop := range schema.Properties {
		fieldName := getGolangName(propKey)
		// calculate sub-schema name here, may not actually be used depending on type of schema!
		subSchemaName := g.getSchemaName(fieldName, prop)
		fieldType, err := g.processSchema(subSchemaName, prop)
		if err != nil {
			return "", err
		}
		f := Field{
			Name:        fieldName,
			JSONName:    propKey,
			Type:        fieldType,
			Required:    isRequired(schema.Required, propKey, prop),
			Description: prop.Description,
			Format:      prop.Format,
			EnumValues:  prop.EnumValues,
		}
		if f.Type == "string" && f.Format == "date-time" {
			strct.importTypes = append(strct.importTypes, "time")
		} else if f.Format == "raw" {
			strct.importTypes = append(strct.importTypes, "encoding/json")
		}
		if f.Required && g.config.EnforceRequiredInMarshallers {
			strct.GenerateCode = true
		}
		strct.Fields[f.Name] = f
	}
	// additionalProperties with typed sub-schema
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.AdditionalPropertiesBool == nil {
		ap := (*Schema)(schema.AdditionalProperties)
		apName := g.getSchemaName("", ap)
		subTyp, err := g.processSchema(apName, ap)
		if err != nil {
			return "", err
		}
		mapTyp := "map[string]" + strings.TrimPrefix(subTyp, "*")
		// If this object is inline property for another object, and only contains additional properties, we can
		// collapse the structure down to a map.
		//
		// If this object is a definition and only contains additional properties, we can't do that or we end up with
		// no struct
		isDefinitionObject := strings.HasPrefix(schema.PathElement, "definitions")
		if len(schema.Properties) == 0 && !isDefinitionObject {
			// since there are no regular properties, we don't need to emit a struct for this object - return the
			// additionalProperties map type.
			return mapTyp, nil
		}

		// if it has no internals, then just alias it to a map
		if len(schema.Properties) == 0 {
			g.Aliases[name] = Field{Name: name, Type: mapTyp}
			return getPrimitiveTypeName("object", name, false)
		}

		// this struct will have both regular and additional properties
		f := Field{
			Name:        "AdditionalProperties",
			JSONName:    "-",
			Type:        mapTyp,
			Required:    false,
			Description: "",
		}
		strct.Fields[f.Name] = f
		// setting this will cause marshal code to be emitted in Output()
		strct.GenerateCode = true
		strct.AdditionalType = subTyp
	}
	// additionalProperties as either true (everything) or false (nothing)
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.AdditionalPropertiesBool != nil {
		if *schema.AdditionalProperties.AdditionalPropertiesBool {
			// everything is valid additional
			subTyp := "map[string]interface{}"
			f := Field{
				Name:        "AdditionalProperties",
				JSONName:    "-",
				Type:        subTyp,
				Required:    false,
				Description: "",
			}
			strct.Fields[f.Name] = f
			// setting this will cause marshal code to be emitted in Output()
			strct.GenerateCode = true
			strct.AdditionalType = "interface{}"
		} else {
			// nothing
			strct.GenerateCode = true
			strct.AdditionalType = "false"
		}
	}
	g.Structs[strct.Name] = strct
	// objects are always a pointer
	return getPrimitiveTypeName("object", name, true)
}

func isRequired(requiredArr []string, propKey string, schema *Schema) bool {
	if schema != nil && schema.AllowsNull() {
		return false
	}

	return contains(requiredArr, propKey)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getPrimitiveTypeName(schemaType string, subType string, pointer bool) (name string, err error) {
	switch schemaType {
	case "array":
		if subType == "" {
			return "error_creating_array", errors.New("can't create an array of an empty subtype")
		}
		if subType == "int" {
			// subType = "int64"
			subType = "float64"
		}
		return "[]" + strings.TrimPrefix(subType, "*"), nil
	case "boolean":
		return "bool", nil
	case "integer":
		// return "int64", nil
		return "float64", nil
	case "number":
		return "float64", nil
	case "null":
		return "nil", nil
	case "object":
		if subType == "" {
			return "error_creating_object", errors.New("can't create an object of an empty subtype")
		}
		if pointer {
			return "*" + subType, nil
		}
		return subType, nil
	case "string":
		return "string", nil
	}

	return "undefined", fmt.Errorf("failed to get a primitive type for schemaType %s and subtype %s",
		schemaType, subType)
}

// return a name for this (sub-)schema.
func (g *Generator) getSchemaName(keyName string, schema *Schema) string {
	if keyName != "" {
		return getGolangName(keyName)
	}
	if schema.Parent == nil {
		return g.getRootElementName(schema)
	}
	if len(schema.Title) > 0 {
		return getGolangName(schema.Title)
	}
	if schema.JSONKey != "" {
		return getGolangName(schema.JSONKey)
	}
	if schema.Parent != nil && schema.Parent.JSONKey != "" {
		return getGolangName(schema.Parent.JSONKey + "Item")
	}
	return g.getAnonymousType()
}

func (g *Generator) getRootElementName(schema *Schema) string {
	// return the hard coded root they wanted
	if g.config.RootElement != "" {
		return g.config.RootElement
	}

	// if this is an aws event, assume they want AwsEvent
	if schema.AwsDetailType != "" || schema.AwsSource != "" {
		return "AwsEvent"
	}

	// default
	return "Root"
}

func (g *Generator) getAnonymousType() string {
	g.anonCount++
	return fmt.Sprintf("Anonymous%d", g.anonCount)
}

// getGolangName strips invalid characters out of golang struct or field names.
func getGolangName(s string) string {
	buf := bytes.NewBuffer([]byte{})
	for i, v := range splitOnAll(s, isNotAGoNameCharacter) {
		if i == 0 && strings.IndexAny(v, "0123456789") == 0 {
			// Go types are not allowed to start with a number, lets prefix with an underscore.
			buf.WriteRune('_')
		}
		buf.WriteString(capitaliseFirstLetter(v))
	}
	return buf.String()
}

func splitOnAll(s string, shouldSplit func(r rune) bool) []string {
	rv := []string{}
	buf := bytes.NewBuffer([]byte{})
	for _, c := range s {
		if shouldSplit(c) {
			rv = append(rv, buf.String())
			buf.Reset()
		} else {
			buf.WriteRune(c)
		}
	}
	if buf.Len() > 0 {
		rv = append(rv, buf.String())
	}
	return rv
}

func isNotAGoNameCharacter(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return true
}

func capitaliseFirstLetter(s string) string {
	if s == "" {
		return s
	}
	prefix := s[0:1]
	suffix := s[1:]
	return strings.ToUpper(prefix) + suffix
}

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, 0, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us
}
