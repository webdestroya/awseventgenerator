//go:build generate
// +build generate

package main

import (
	_ "embed"
	"log"
	"os"
	"path"
	"strings"

	"github.com/webdestroya/awseventgenerator"
	"github.com/webdestroya/awseventgenerator/internal/testutil/testwriter"
)

var (
	jsonPath   = "../testdata"
	goCodePath = "!SET_BELOW!"
	testGenDir = "!SET_BELOW!"
)

const (
	// true = actually write/delete test files
	// false = pretend
	allowFSWriting = true
)

var testConfigs = map[string]awseventgenerator.Config{
	"normal": {
		GenerateEnums: true,
	},
	"alwaysptr": {
		GenerateEnums:    true,
		AlwaysPointerize: true,
	},
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("[TestPackage Generator] ")

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	goCodePath = workDir

	testGenDir = path.Join(goCodePath, "testsuite_gen")

	cleanUpExisting()
	generateTestPackages()
}

func cleanUpExisting() {
	if !allowFSWriting {
		return
	}

	folders := make([]string, 0, len(testConfigs))

	for k := range testConfigs {
		folders = append(folders, k+"_gen")
	}

	for _, folder := range folders {
		dirPath := path.Join(goCodePath, folder)

		log.Printf("Deleting Directory: %s", dirPath)
		if err := os.RemoveAll(dirPath); err != nil {
			log.Fatalf("Failed to delete folder: %s %s", dirPath, err)
		}
	}

}

func generateTestPackages() {
	files, err := os.ReadDir(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if path.Ext(file.Name()) != ".json" {
			continue
		}

		// if file.Name() != "aprefnoprop.json" {
		// 	continue
		// }

		for k, v := range testConfigs {
			genTestForConfig(k, file.Name(), v)
		}

	}
}

func genTestForConfig(label, filename string, config awseventgenerator.Config) {
	packageName := strings.TrimSuffix(strings.ToLower(filename), ".json")

	parentFolder := label + "_gen"
	// testFile := packageName + "_" + label + "_test.go"

	jsonFile := path.Join(jsonPath, filename)
	log.Printf("Writing %s::%s for %s", label, packageName, filename)

	folderName := packageName + "_gen"

	destFile := path.Join(goCodePath, parentFolder, folderName, "generated.go")
	testFile := path.Join(goCodePath, parentFolder, folderName, "generated_test.go")

	config.PackageName = packageName

	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("failed to read json file: %s", err)
	}

	schema, err := awseventgenerator.Parse(string(jsonData), nil)
	if err != nil {
		log.Fatalf("failed to parse schema: %s", err)
	}

	data, err := awseventgenerator.GenerateFromSchema(schema, &config)
	if err != nil {
		log.Fatalf("Failure: %s => %s", filename, err)
	}

	twriter := testwriter.NewTestWriter()
	if err := twriter.Add(data, label, packageName, path.Join(parentFolder, folderName), schema); err != nil {
		log.Fatalf("Failed to add to testwriter: %s %s", packageName, err)
	}

	writeFile(destFile, data)

	testBytes, err := twriter.Generate()
	if err != nil {
		log.Println(string(testBytes))
		log.Fatalf("Failed to generate test: %s", err)
	}

	writeFile(testFile, testBytes)
}

func writeFile(filename string, data []byte) {

	if !allowFSWriting {
		log.Printf("Would write file %s", filename)
		return
	}

	destDir := path.Dir(filename)

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("Could not make directories for: %s %s", destDir, err)
		}
	}

	if err := os.WriteFile(filename, data, 0o600); err != nil {
		log.Fatalf("Could not write file: %s %s", filename, err)
	}

}
