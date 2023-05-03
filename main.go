package awseventgenerator

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"path"
)

func GenerateFromSchemaFile(filename string, config *Config) ([]byte, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	abPath, err := abs(filename)
	if err != nil {
		return nil, errors.New("failed to normalise input path with error " + err.Error())
	}

	fileURI := &url.URL{
		Scheme: "file",
		Path:   abPath,
	}

	schema, err := parseWithSchemaKeyRequired(string(data), fileURI, false)
	if err != nil {
		return nil, err
	}

	return GenerateFromSchema(schema, config)
}

func GenerateFromSchemaString(data string, config *Config) ([]byte, error) {
	fileURI := &url.URL{
		Scheme: "file",
		Path:   "stringdata.json",
	}

	schema, err := parseWithSchemaKeyRequired(data, fileURI, false)
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

func abs(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, name), err
}
