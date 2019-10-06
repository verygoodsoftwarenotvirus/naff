package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/client/v1/http"
	cmdv1server "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/cmd/server/v1"
	twofactor "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates/cmd/tools/two_factor"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {

	todoDataTypes := []models.DataType{
		{
			Name: models.Name{
				Singular:                "Item",
				Plural:                  "Items",
				RouteName:               "items",
				PluralRouteName:         "item",
				UnexportedVarName:       "item",
				PluralUnexportedVarName: "items",
			},
		},
	}

	client.RenderPackage(todoDataTypes)
	cmdv1server.RenderPackage(todoDataTypes)
	twofactor.RenderPackage(todoDataTypes)

	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http/helpers.go")
	outputPath := strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/example_output", 1)

	dmp := diffmatchpatch.New()
	diffs, err := checkFiles(dmp, sourcePath, outputPath)
	if err != nil {
		log.Fatal(err)
	}

	if diffs != nil {
		x := dmp.DiffPrettyText(diffs)
		println(x)
	}
}

func checkFiles(dmp *diffmatchpatch.DiffMatchPatch, path1, path2 string) ([]diffmatchpatch.Diff, error) {
	file1, err := ioutil.ReadFile(path1)
	if err != nil {
		return nil, err
	}
	file2, err := ioutil.ReadFile(path2)
	if err != nil {
		return nil, err
	}

	diffs := dmp.DiffMain(string(file1), string(file2), false)

	if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
		return nil, nil
	}

	return diffs, nil
}
