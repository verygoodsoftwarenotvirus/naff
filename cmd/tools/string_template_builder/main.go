package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	allPackages := []string{
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/fakes",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/client/httpclient/requests",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/client/httpclient",
	}

	for _, pkg := range allPackages {
		err := doTheThingForPackage(pkg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doTheThingForPackage(pkgPath string) error {
	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", pkgPath)
	outputPath := strings.Replace(strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/templates/experimental", 1), "_test.go", "test_.go", 1)

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fp := filepath.Join(sourcePath, file.Name())
			op := strings.TrimPrefix(
				strings.Replace(filepath.Join(outputPath, file.Name()), filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff/"), "", 1),
				"/",
			)
			op = strings.Replace(op, ".go", ".gotpl", 1)
			op = strings.Replace(op, "/internal/", "/_internal_/", 1)

			// i don't care
			_ = os.MkdirAll(filepath.Dir(op), 0777)

			var fileBytes []byte
			fileBytes, err = ioutil.ReadFile(fp)
			if err != nil {
				return fmt.Errorf("error reading input file: %w", err)
			}

			fileContents := regexp.MustCompile(`\t(\w+\s)?"gitlab\.com\/verygoodsoftwarenotvirus\/todo\/([\w\/]+)"`).ReplaceAllString(string(fileBytes), `	$1{{ projectImport "$2" }}`)
			fileContents = regexp.MustCompile(`(?i)todo`).ReplaceAllString(fileContents, `{{ projectName }}`)

			if err = ioutil.WriteFile(op, []byte(fileContents), 0644); err != nil {
				return fmt.Errorf("error writing output file: %w", err)
			}
		}
	}

	return nil
}
