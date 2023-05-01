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
	goCodePath = ""
)

const (
	// true = actually write/delete test files
	// false = pretend
	allowFSWriting = true
)

//go:embed helpers.go.tmpl
var helpersTmpl []byte

func main() {
	log.SetPrefix("[TestPackage Generator] ")

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	goCodePath = workDir

	// cleanUpExisting()
	generateTestPackages()
}

func cleanUpExisting() {
	if !allowFSWriting {
		return
	}
	files, err := os.ReadDir(goCodePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), "_gen") {
			continue
		}

		dirPath := path.Join(goCodePath, file.Name())

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

	testGenDir := path.Join(goCodePath, "testsuite")
	if _, err := os.Stat(testGenDir); os.IsNotExist(err) {
		if err := os.MkdirAll(testGenDir, os.ModePerm); err != nil {
			log.Fatalf("Could not make directories for: %s %s", testGenDir, err)
		}
	}
	if err := os.WriteFile(path.Join(testGenDir, "utils_test.go"), helpersTmpl, 0o600); err != nil {
		log.Fatalf("Could not write test helpers: %s", err)
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

		packageName := strings.TrimSuffix(strings.ToLower(file.Name()), ".json")
		folderName := packageName + "_gen"
		jsonFile := path.Join(jsonPath, file.Name())

		destFile := path.Join(goCodePath, folderName, "generated.go")
		log.Printf("Writing %s for %s", folderName, file.Name())

		config := &awseventgenerator.Config{
			PackageName:             packageName,
			RootElement:             "Root",
			GenerateEnums:           true,
			GenerateEnumValueMethod: true,
		}

		data, err := awseventgenerator.GenerateFromSchemaFile(jsonFile, config)
		if err != nil {
			log.Fatalf("Failure: %s => %s", file.Name(), err)
		}

		destDir := path.Dir(destFile)

		twriter := testwriter.NewTestWriter()
		if err := twriter.Add(data, packageName, folderName); err != nil {
			log.Fatalf("Failed to add to testwriter: %s %s", packageName, err)
		}

		if allowFSWriting {
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
					log.Fatalf("Could not make directories for: %s %s", destDir, err)
				}
			}

			if err := os.WriteFile(destFile, data, 0o600); err != nil {
				log.Fatalf("Could not write file: %s %s", destFile, err)
			}

			testBytes, err := twriter.Generate()
			if err != nil {
				log.Fatalf("Failed to generate test: %s", err)
			}

			if err := os.WriteFile(path.Join(testGenDir, packageName+"_test.go"), testBytes, 0o600); err != nil {
				log.Fatalf("Could not write generated test: %s", err)
			}
		}

	}

	// fmt.Println("TEST", string(testBytes))

}
