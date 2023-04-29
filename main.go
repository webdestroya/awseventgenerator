package awseventgenerator

func GenerateFromSchemaFile(filename string, config *Config) ([]byte, error) {
	return GenerateFromSchema(nil, config)
}

func GenerateFromSchemaString(filename string, config *Config) ([]byte, error) {
	return GenerateFromSchema(nil, config)
}

func GenerateFromSchema(schema *Schema, config *Config) ([]byte, error) {
	return nil, nil
}
