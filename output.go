package awseventgenerator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	skipCodeGen = false
)

var (
	titleCaser    = cases.Title(language.AmericanEnglish)
	wordSepRegexp = regexp.MustCompile(`[-_]`)
)

type outputter struct {
	g *Generator
}

// Output generates code and writes to w.
func Output(w io.Writer, g *Generator) error {

	var outputBuf bytes.Buffer

	structs := g.Structs
	aliases := g.Aliases

	var pkg string

	switch pVal := g.config.PackageName.(type) {
	case string:
		pkg = pVal
	case PackageNameFunc:
		pkg = pVal(g.schemas[0])
	case nil:
		pkg = DefaultPackageNameFunc(g.schemas[0])
	default:
		return fmt.Errorf("invalid PackageName config value: %T %v", pVal, pVal)
	}

	fmt.Fprintln(&outputBuf, "// Code generated by github.com/webdestroya/awseventgenerator. DO NOT EDIT.")
	fmt.Fprintln(&outputBuf)
	fmt.Fprintf(&outputBuf, "package %v\n", cleanPackageName(pkg))

	// write all the code into a buffer, compiler functions will return list of imports
	// write list of imports into main output stream, followed by the code
	codeBuf := new(bytes.Buffer)
	imports := make(map[string]bool)

	topLevel := collectTopLevelStruct(structs)
	_ = topLevel

	for _, k := range getSortedKeys(structs) {
		s := structs[k]
		s.finalize(g)
		if !skipCodeGen {
			if s.GenerateCode {
				emitMarshalCode(codeBuf, s, imports, g.config)
				emitUnmarshalCode(codeBuf, s, imports, g.config)
			}
		} else {
			imports["encoding/json"] = true
		}
		for _, str := range s.importTypes {
			imports[str] = true
		}
	}

	if len(imports) > 0 {
		fmt.Fprintf(&outputBuf, "\nimport (\n")
		for _, k := range getSortedKeys(imports) {
			fmt.Fprintf(&outputBuf, "    \"%s\"\n", k)
		}
		fmt.Fprintf(&outputBuf, ")\n")
	}

	if len(g.Constants) > 0 {
		fmt.Fprintf(&outputBuf, "\nconst (\n")
		for _, k := range getSortedKeys(g.Constants) {
			vraw := g.Constants[k]
			fmt.Fprintf(&outputBuf, "  %s = ", k)
			switch v := vraw.(type) {
			case int, float64:
				fmt.Fprintf(&outputBuf, "%v\n", v)
			case bool:
				fmt.Fprintf(&outputBuf, "%t\n", v)
			default:
				fmt.Fprintf(&outputBuf, "`%v`\n", v)
			}
		}
		fmt.Fprintf(&outputBuf, ")\n")
	}

	for _, k := range getSortedKeys(aliases) {
		a := aliases[k]

		fmt.Fprintln(&outputBuf, "")
		// fmt.Fprintf(&outputBuf, "// %s\n", a.Name)
		fmt.Fprintf(&outputBuf, "type %s = %s\n", a.Name, a.Type)
	}

	// ENUMS
	for _, k := range getSortedKeys(structs) {
		s := structs[k]

		if !s.IsEnum {
			continue
		}

		fmt.Fprintln(&outputBuf, "")
		outputNameAndDescriptionComment(s.Name, s.Description, &outputBuf)
		if err := emitEnum(&outputBuf, g, s); err != nil {
			return err
		}
	}

	for _, k := range getSortedKeys(structs) {
		s := structs[k]

		if s.IsEnum {
			continue
		}

		fmt.Fprintln(&outputBuf, "")
		outputNameAndDescriptionComment(s.Name, s.Description, &outputBuf)

		// if len(s.Fields) == 0 {
		// 	if s.AdditionalType == "" {
		// 		fmt.Fprintf(&outputBuf, "type %s = interface{}\n", s.Name)
		// 	} else {
		// 		fmt.Fprintf(&outputBuf, "type %s = interface{} // Additional Type: %s\n", s.Name, s.AdditionalType)
		// 	}
		// 	continue
		// }
		if s.AliasType != "" {
			fmt.Fprintf(&outputBuf, "type %s = %s\n", s.Name, s.AliasType)
			continue
		}

		fmt.Fprintf(&outputBuf, "type %s struct {\n", s.Name)

		for _, fieldKey := range getSortedKeys(s.Fields) {
			f := s.Fields[fieldKey]

			if !f.finalized {
				panic("NOT FINIALIZED")
			}

			// Only apply omitempty if the field is not required.
			omitempty := ",omitempty"
			if f.Required && !g.config.AlwaysPointerize {
				omitempty = ""
			}

			if f.Description != "" {
				outputFieldDescriptionComment(f.Description, &outputBuf)
			}

			if !g.config.GenerateEnums && len(f.EnumValues) > 0 {
				fmt.Fprintf(&outputBuf, "\n  // Values: %s\n", strings.Join(f.EnumValues, ", "))
			}

			jsonName := f.JSONName
			ftype := g.finalFieldTypes[s.Name+":"+f.Name]
			fmt.Fprintf(&outputBuf, "  %s %s `json:\"%s%s\"`\n", f.Name, ftype, jsonName, omitempty)
		}

		fmt.Fprintln(&outputBuf, "}")
	}

	// write code after structs for clarity
	outputBuf.Write(codeBuf.Bytes())

	if g.config.NoFormatCode {
		_, err := outputBuf.WriteTo(w)
		return err

	} else {

		formattedBytes, err := format.Source(outputBuf.Bytes())
		if err != nil {
			return err
		}

		_, err = w.Write(formattedBytes)
		return err
	}
}

func emitEnum(w io.Writer, g *Generator, s Struct) error {

	fmt.Fprintf(w, "type %s string\n", s.Name)
	fmt.Fprintln(w, `const (`)
	for _, val := range s.EnumValues {
		fmt.Fprintf(w, "  %s%s %s = \"%s\"\n", s.Name, enumifyValue(val), s.Name, val)
	}
	fmt.Fprintln(w, `)`)

	fmt.Fprintf(w, "func (%s) Values() []%s {\n", s.Name, s.Name)
	fmt.Fprintf(w, "  return []%s{\n", s.Name)
	for _, val := range s.EnumValues {
		fmt.Fprintf(w, "  \"%s\",\n", val)
	}
	fmt.Fprintln(w, "  }")
	fmt.Fprintln(w, "}")

	return nil
}

func emitMarshalCode(w io.Writer, s Struct, imports map[string]bool, config *Config) {
	fmt.Fprintf(w, "func (strct *%s) MarshalJSON() ([]byte, error) {", s.Name)
	fmt.Fprintln(w, "  data := make(map[string]interface{})")

	if len(s.Fields) > 0 {
		// Marshal all the defined fields
		for _, fieldKey := range getSortedKeys(s.Fields) {
			f := s.Fields[fieldKey]
			if f.JSONName == "-" {
				continue
			}

			if f.Required {
				fmt.Fprintf(w, "  data[\"%s\"] = strct.%s\n", f.JSONName, f.Name)
			} else {
				fmt.Fprintf(w, "  if strct.%s != nil {\n", f.Name)
				fmt.Fprintf(w, "  data[\"%s\"] = strct.%s\n", f.JSONName, f.Name)
				fmt.Fprintln(w, "  }")
			}
		}
	}

	if s.AdditionalType != "" && s.AdditionalType != "false" {
		fmt.Fprintln(w, `    for k, v := range strct.AdditionalProperties {`)
		fmt.Fprintln(w, `    data[k] = v`)
		fmt.Fprintln(w, `    }`)

	}

	fmt.Fprintln(w, "  return json.Marshal(data)")
	fmt.Fprintln(w, "}")
}

func emitUnmarshalCode(w io.Writer, s Struct, imports map[string]bool, config *Config) {
	imports["encoding/json"] = true
	// unmarshal code
	fmt.Fprintf(w, `
func (strct *%s) UnmarshalJSON(b []byte) error {
`, s.Name)
	// setup required bools
	for _, fieldKey := range getSortedKeys(s.Fields) {
		f := s.Fields[fieldKey]
		if f.Required && config.EnforceRequiredInMarshallers {
			fmt.Fprintf(w, "    %sReceived := false\n", f.JSONName)
		}
	}
	// setup initial unmarshal
	fmt.Fprintf(w, `    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }`)

	// figure out if we need the "v" output of the range keyword
	needVal := "_"
	if len(s.Fields) > 0 || s.AdditionalType != "false" {
		needVal = "v"
	}
	// start the loop
	fmt.Fprintf(w, `
    // parse all the defined properties
    for k, %s := range jsonMap {
        switch k {
`, needVal)
	// handle defined properties
	for _, fieldKey := range getSortedKeys(s.Fields) {
		f := s.Fields[fieldKey]
		if f.JSONName == "-" {
			continue
		}
		fmt.Fprintf(w, `        case "%s":
            if err := json.Unmarshal([]byte(v), &strct.%s); err != nil {
                return err
             }
`, f.JSONName, f.Name)
		if f.Required && config.EnforceRequiredInMarshallers {
			fmt.Fprintf(w, "            %sReceived = true\n", f.JSONName)
		}
	}

	// handle additional property
	if s.AdditionalType != "" {
		cleanedAddtlType := strings.TrimPrefix(s.AdditionalType, "*")
		if s.AdditionalType == "false" {
			// all unknown properties are not allowed
			imports["fmt"] = true
			fmt.Fprintf(w, `        default:
            return fmt.Errorf("additional property not allowed: \"" + k + "\"")
`)
		} else {
			fmt.Fprintf(w, `        default:
            // an additional "%s" value
            var additionalValue %s
            if err := json.Unmarshal([]byte(v), &additionalValue); err != nil {
                return err // invalid additionalProperty
            }
            if strct.AdditionalProperties == nil {
                strct.AdditionalProperties = make(map[string]%s, 0)
            }
            strct.AdditionalProperties[k]= additionalValue
`, cleanedAddtlType, cleanedAddtlType, cleanedAddtlType)
		}
	}
	fmt.Fprintf(w, "        }\n") // switch
	fmt.Fprintf(w, "    }\n")     // for

	// check all Required fields were received
	if config.EnforceRequiredInMarshallers {
		for _, fieldKey := range getSortedKeys(s.Fields) {
			f := s.Fields[fieldKey]
			if f.Required {
				imports["errors"] = true
				fmt.Fprintf(w, `    // check if %s (a required property) was received
    if !%sReceived {
        return errors.New("\"%s\" is required but was not present")
    }
`, f.JSONName, f.JSONName, f.JSONName)
			}
		}
	}

	fmt.Fprintf(w, "    return nil\n")
	fmt.Fprintf(w, "}\n") // UnmarshalJSON
}

func outputNameAndDescriptionComment(name, description string, w io.Writer) {
	if description == "" {
		return
	}

	if !strings.Contains(description, "\n") {
		fmt.Fprintf(w, "// %s\n", description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "// %s\n", strings.Join(dl, "\n// "))
}

func outputFieldDescriptionComment(description string, w io.Writer) {
	if !strings.Contains(description, "\n") {
		fmt.Fprintf(w, "\n  // %s\n", description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "\n  // %s\n", strings.Join(dl, "\n  // "))
}

func enumifyValue(v string) string {
	v = wordSepRegexp.ReplaceAllString(v, " ")
	v = titleCaser.String(v)

	return strings.ReplaceAll(v, " ", "")
}

func cleanPackageName(pkg string) string {
	pkg = strings.ReplaceAll(pkg, " ", "")
	pkg = strings.ReplaceAll(pkg, ".", "")
	pkg = strings.ReplaceAll(pkg, "-", "")
	return pkg
}

func getSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func collectTopLevelStruct(structs map[string]Struct) map[string]Struct {
	m := make(map[string]Struct)

	for sk, s := range structs {
		m[sk] = s
	}

	for _, s := range structs {
		for _, f := range s.Fields {
			sname := strings.Trim(f.Type, "[]*")
			delete(m, sname)
		}
	}
	return m
}
