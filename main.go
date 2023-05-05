package awseventgenerator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
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

// data =
// string = file (unless it starts with {)
// []byte = interpreted as raw schema
// *Schema
// io.Reader
func Generate(from any, config *Config) ([]byte, error) {
	switch val := from.(type) {
	case []byte:
		return GenerateFromSchemaString(string(val), config)

	case string:
		// if it smells like json, just assume it is
		if strings.ContainsAny(val, `{},"[]`) {
			return GenerateFromSchemaString(val, config)
		}

		return GenerateFromSchemaFile(val, config)
	case io.Reader:
		filedata, err := io.ReadAll(val)
		if err != nil {
			return nil, err
		}
		return GenerateFromSchemaString(string(filedata), config)
	case *Schema:
		return GenerateFromSchema(val, config)
	default:
		return nil, fmt.Errorf("invalid source type: %T", from)
	}
}

func GenerateAndExport(data any, destinationFile string, config *Config) error {

	output, err := Generate(data, config)
	if err != nil {
		return err
	}

	destDir := path.Dir(destinationFile)

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return fmt.Errorf("Could not make directories for: %s %w", destDir, err)
		}
	}

	if err := os.WriteFile(destinationFile, output, 0o600); err != nil {
		return fmt.Errorf("Could not write file: %s %w", destinationFile, err)
	}

	return nil
}

func abs(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, name), err
}
