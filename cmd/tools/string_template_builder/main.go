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
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/server",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/config_gen",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/data_scaffolder",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/encoded_qr_code_generator",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/index_initializer",
		"gitlab.com/verygoodsoftwarenotvirus/todo/cmd/tools/template_gen",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/authentication",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/authorization",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/build/server",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/capitalism",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/capitalism/stripe",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/config",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/viper",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querier",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mariadb",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/postgres",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/querybuilding/sqlite",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/events",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/keys",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/logging",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/tracing",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/panicking",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/random",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/chi",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/search",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/bleve",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/secrets",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/server",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/accounts",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/admin",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/apiclients",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/audit",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/frontend",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/items",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/users",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/webhooks",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/storage",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/images",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/client/httpclient",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/client/httpclient/requests",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/converters",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/fakes",
		"gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/frontend",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/integration",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/load",
		"gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils",
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
	outputPath := strings.Replace(strings.Replace(sourcePath, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/templates", 1), "_test.go", "test_.go", 1)

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

			// i don't care
			_ = os.MkdirAll(filepath.Dir(op), 0777)

			var fileBytes []byte
			fileBytes, err = ioutil.ReadFile(fp)
			if err != nil {
				return fmt.Errorf("error reading input file: %w", err)
			}

			fileContents := regexp.MustCompile(`\t(\w+\s)?"gitlab\.com\/verygoodsoftwarenotvirus\/todo\/([\w\/]+)`).ReplaceAllString(string(fileBytes), `	$1"{{ projectImport "$2" }}`)
			fileContents = regexp.MustCompile(`(?i)todo`).ReplaceAllString(fileContents, `{{ projectName }}`)

			if err = ioutil.WriteFile(op, []byte(fileContents), 0644); err != nil {
				return fmt.Errorf("error writing output file: %w", err)
			}
		}
	}

	return nil
}
