package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	outputRepoVariableTemplate   = "{{ .OutputRepository }}"
	sourceRepositoryPath         = "gitlab.com/verygoodsoftwarenotvirus/todo"
	iterablesImportsTemplateCode = `{{ range $x, $import := .IterableServicesImports }} "$import" {{ end }}`
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gopath := os.Getenv("GOPATH")
	sourcePath := filepath.Join(gopath, "src", sourceRepositoryPath)

	skipDirectories := map[string]bool{
		".idea":                    true,
		".git":                     true,
		"artifacts":                true,
		"vendor":                   true,
		"frontend/v1/node_modules": true,
	}

	skipFiles := map[string]bool{
		"frontend/v1/public/bundle.css":     true,
		"frontend/v1/public/bundle.css.map": true,
		"frontend/v1/public/bundle.js":      true,
		"frontend/v1/public/bundle.js.map":  true,
		"cmd/server/v1/wire_gen.go":         true,
	}

	iterableDirectories := map[string]bool{
		"frontend/v1/src/pages/items": true,
		"services/v1/items":           true,
	}

	iterableFiles := map[string]bool{
		"client/v1/http/items.go":      true,
		"client/v1/http/items_test.go": true,

		"database/v1/client/items.go":      true,
		"database/v1/client/items_test.go": true,

		"database/v1/queriers/postgres/items.go":      true,
		"database/v1/queriers/postgres/items_test.go": true,

		"frontend/v1/src/pages/items/Create.svelte": true,
		"frontend/v1/src/pages/items/List.svelte":   true,
		"frontend/v1/src/pages/items/Read.svelte":   true,

		"models/v1/item.go":      true,
		"models/v1/item_test.go": true,

		"models/v1/mock/mock_item_data_manager.go": true,
		"models/v1/mock/mock_item_data_server.go":  true,

		"services/v1/items/http_routes.go":        true,
		"services/v1/items/http_routes_test.go":   true,
		"services/v1/items/items_service.go":      true,
		"services/v1/items/items_service_test.go": true,
		"services/v1/items/middleware.go":         true,
		"services/v1/items/middleware_test.go":    true,
		"services/v1/items/wire.go":               true,

		"tests/v1/testutil/rand/model/items.go": true,
		"tests/v1/integration/items_test.go":    true,
		"tests/v1/load/items.go":                true,
	}

	// files that require exceptional handling
	specialSnowflakes := map[string]func(string) string{
		"server/v1/http/routes.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			lines := strings.Split(in, "\n")
			output := []string{}
			var removingLines bool
			removedLines := []string{}

			for _, line := range lines {
				trimLine := strings.TrimSpace(line)
				if trimLine == "// Items" {
					removingLines = true
					removedLines = append(removedLines, line)
				} else if removingLines && trimLine != "})" {
					removedLines = append(removedLines, line)
					continue
				} else if removingLines && trimLine == "})" {
					removedLines = append(removedLines, line)
					removingLines = false
				} else {
					output = append(output, line)
				}
			}

			return strings.Join(output, "\n")
		},
		"deploy/grafana/dashboards/dashboard.json": func(in string) string {
			itemCounterChunk := `
				{
					"expr": "todo_server_items_count",
					"format": "time_series",
					"instant": false,
					"intervalFactor": 1,
					"legendFormat": "items",
					"refId": "A"
				},`

			return strings.Replace(in, itemCounterChunk, "", 1)
		},
		"cmd/server/v1/wire.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(
				in, `		items.Providers,`, "", 1)

			return in
		},
		"database/v1/database.go": func(in string) string {
			in = strings.Replace(in, `
		models.ItemDataManager`, "", 1)

			return in
		},
		"database/v1/queriers/postgres/migrations.go": func(in string) string {
			in = strings.Replace(in, `		{
			Version:     4,
			Description: "create items table",
			Script: `, "", 1)

			in = strings.Replace(in, `
			CREATE TABLE IF NOT EXISTS items (
				"id" bigserial NOT NULL PRIMARY KEY,
				"name" text NOT NULL,
				"details" text NOT NULL DEFAULT '',
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				"belongs_to" bigint NOT NULL,
				FOREIGN KEY ("belongs_to") REFERENCES "users"("id")
			);`, "", 1)

			in = strings.Replace(in, "``,\n\t\t},", "", 1)

			return in
		},
		"server/v1/http/server.go": func(in string) string {
			in = strings.Replace(in, "		itemsService         models.ItemDataServer", "\n", 1)
			in = strings.Replace(in, "	itemsService models.ItemDataServer,", "", 1)
			in = strings.Replace(in, "		itemsService:         itemsService,", "", 1)

			return in
		},
		"server/v1/http/server_test.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)
			in = strings.Replace(in, "		itemsService:         &mmodels.ItemDataServer{},", "", 1)
			in = strings.Replace(in, "			&items.Service{},", "", 1)

			return in
		},
		"server/v1/http/wire_param_fetchers.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(in, "		ProvideItemIDFetcher,", "", 1)

			in = strings.Replace(in, `// ProvideItemIDFetcher provides an ItemIDFetcher
func ProvideItemIDFetcher(logger logging.Logger) items.ItemIDFetcher {
	return buildChiItemIDFetcher(logger)
}`, "", 1)

			in = strings.Replace(in, `

// chiItemIDFetcher fetches a ItemID from a request routed by chi.
func buildChiItemIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, items.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ItemID from request")
		}
		return u
	}
}`, "", 1)

			in = strings.Replace(in, `
// ProvideUserIDFetcher provides a UserIDFetcher
func ProvideUserIDFetcher() items.UserIDFetcher {
	return UserIDFetcher
}`, "", 1)

			in = strings.Replace(in, `		ProvideUserIDFetcher,`, "", 1)

			return in
		},
		"server/v1/http/wire_param_fetchers_test.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(in, `
func TestProvideUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUserIDFetcher()
	})
}
`, "", 1)

			in = strings.Replace(in, `ovideItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemIDFetcher(noop.ProvideNoopLogger())
	})
}`, "", 1)

			in = strings.Replace(in, `
func Test_buildChiItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildChiItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{items.URIParamKey},
				Values: []string{fmt.Sprintf("%d", expected)},
			},
		}))

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildChiItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{items.URIParamKey},
				Values: []string{"expected"},
			},
		}))

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}`, "", 1)

			return in
		},
		"services/v1/frontend/http_routes.go": func(in string) string {
			in = strings.Replace(in, `
	// itemsFrontendPathRegex matches URLs against our frontend router's specification for specific item routes
	itemsFrontendPathRegex = regexp.MustCompile(`, "", 1)

			in = strings.Replace(in, "`/items/\\d+`)", "", 1)
			in = strings.Replace(in, `			"/items",`, "", 1)
			in = strings.Replace(in, `

			"/items/new",`, "", 1)
			in = strings.Replace(in, `
		if itemsFrontendPathRegex.MatchString(req.URL.Path) {
			req.URL.Path = "/"
		}`, "", 1)

			return in
		},
		"services/v1/frontend/http_routes_test.go": func(in string) string {
			in = strings.Replace(in, `
	T.Run("with frontend routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/login"

		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})`, "", 1)

			in = strings.Replace(in, `

	T.Run("with frontend items routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/items/9"

		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
`, "", 1)

			return in
		},
		"database/v1/database_mock.go": func(in string) string {
			in = strings.Replace(in, `
		ItemDataManager:         &mmodels.ItemDataManager{},`, "", 1)
			in = strings.Replace(in, `
	*mmodels.ItemDataManager`, "", 1)

			return in
		},
		"frontend/v1/src/App.svelte": func(in string) string {
			in = strings.Replace(in, `

  // Items routes
  import ReadItem from "./pages/items/Read.svelte";
  import CreateItem from "./pages/items/Create.svelte";
  import Items from "./pages/items/List.svelte";`, "", 1)

			in = strings.Replace(in, `
    <Link to="items">Items</Link>
    <Link to="items/new">Create Item</Link>`, "", 1)

			in = strings.Replace(in, `
    <Route path="items" component={Items} />
    <Route path="items/:id" component={ReadItem} />
    <Route path="items/new" component={CreateItem} />`, "", 1)

			return in
		},
		"tests/v1/integration/auth_test.go": func(in string) string {

			in = strings.Replace(in, `
	test.Run("should only allow users to see their own content", func(t *testing.T) {`, `
	test.Run("should only allow users to see their own content", func(t *testing.T) {
		// NOTE: this function tests that data is only revealed to folks who have the authority to view it
		// by creating OAuth2 clients. If you have a more meaningful data structure in your service, consider revising`, 1)

			return in
		},
		"tests/v1/load/actions.go": func(in string) string {
			in = strings.Replace(in, `

	for k, v := range buildItemActions(c) {`
		allActions[k] = v
	}`, "", 1)
			return in
		},
	}

	standardReplacers := strings.NewReplacer(
		sourceRepositoryPath, outputRepoVariableTemplate,
	)

	must(os.RemoveAll("template"))
	must(os.MkdirAll("template/base_repository", os.ModePerm))
	must(os.MkdirAll("template/iterables", os.ModePerm))

	// build base repository files
	if err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		outputFilepath := strings.Replace(path, "todo", "naff/template/base_repository", 1)
		relativePath := strings.Replace(path, sourcePath+"/", "", 1)

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			if _, ok := skipDirectories[relativePath]; ok {
				return filepath.SkipDir
			}

			if _, ok := iterableDirectories[relativePath]; ok {
				outputFilepath = strings.Replace(outputFilepath, "base_repository", "iterables", 1)
			}

			return os.MkdirAll(outputFilepath, info.Mode())
		} else {
			if _, ok := skipFiles[relativePath]; ok {
				return nil
			}

			if _, ok := iterableFiles[relativePath]; ok {
				outputFilepath = strings.Replace(outputFilepath, "base_repository", "iterables", 1)
				if err := os.MkdirAll(filepath.Dir(outputFilepath), os.ModePerm); err != nil {
					return err
				}
			}

			// do the thing
			fc, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			file := standardReplacers.Replace(string(fc))

			if replacerFunc, ok := specialSnowflakes[relativePath]; ok {
				if replacerFunc != nil {
					file = replacerFunc(file)
				} else {
					return nil
				}
			}

			e := ioutil.WriteFile(outputFilepath, []byte(file), info.Mode())
			return e
		}
	}); err != nil {
		log.Fatal(err)
	}
}
