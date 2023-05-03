package awseventgenerator

import (
	"bytes"
	"net/url"
)

func GenerateFromSchemaFile(filename string, config *Config) ([]byte, error) {

	schemas, err := readInputFiles([]string{filename}, false)
	if err != nil {
		return nil, err
	}

	return GenerateFromSchema(schemas[0], config)
}

func GenerateFromSchemaString(data string, config *Config) ([]byte, error) {
	fileURI := &url.URL{
		Scheme: "file",
		Path:   "stringdata.json",
	}

	schema, err := Parse(data, fileURI)
	if err != nil {
		return nil, err
	}

	return GenerateFromSchema(schema, config)
}

func GenerateFromSchema(schema *Schema, config *Config) ([]byte, error) {

	g := New(config, schema)
	err := g.Generate()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := Output(&buf, g); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
