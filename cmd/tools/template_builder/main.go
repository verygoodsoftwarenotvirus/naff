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
	iterablesImportsTemplateCode = `{{ range $x, $import := .IterableServicesImports }} 
	"{{ $import }}" 
{{ end }}`
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
		"profile.out":                       true,
	}

	// files that require exceptional handling
	specialSnowflakes := map[string]func(string) string{
		"server/v1/http/routes.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(in, `// Items
			v1Router.Route("/items", func(itemsRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, items.URIParamKey)
				itemsRouter.With(s.itemsService.CreationInputMiddleware).
					Post("/", s.itemsService.CreateHandler) // CreateHandler
				itemsRouter.Get(sr, s.itemsService.ReadHandler) // ReadHandler
				itemsRouter.With(s.itemsService.UpdateInputMiddleware).
					Put(sr, s.itemsService.UpdateHandler) // UpdateHandler
				itemsRouter.Delete(sr, s.itemsService.ArchiveHandler) // ArchiveHandler
				itemsRouter.Get("/", s.itemsService.ListHandler)      // ListHandler
			})`, `{{ range $i, $dt := .DataTypes }}
			// {{ camelcase $dt.Name }}s
			v1Router.Route("/{{ lower $dt.Name }}s", func({{ lower $dt.Name }}sRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, {{ lower $dt.Name }}s.URIParamKey)
				{{ lower $dt.Name }}sRouter.With(s.{{ lower $dt.Name }}sService.CreationInputMiddleware).
					Post("/", s.{{ lower $dt.Name }}sService.CreateHandler) // CreateHandler
				{{ lower $dt.Name }}sRouter.Get(sr, s.{{ lower $dt.Name }}sService.ReadHandler) // ReadHandler
				{{ lower $dt.Name }}sRouter.With(s.{{ lower $dt.Name }}sService.UpdateInputMiddleware).
					Put(sr, s.{{ lower $dt.Name }}sService.UpdateHandler) // UpdateHandler
				{{ lower $dt.Name }}sRouter.Delete(sr, s.{{ lower $dt.Name }}sService.ArchiveHandler) // ArchiveHandler
				{{ lower $dt.Name }}sRouter.Get("/", s.{{ lower $dt.Name }}sService.ListHandler)      // ListHandler
			})
			{{ end }}`, 1)

			return in
		},
		"deploy/grafana/dashboards/dashboard.json": func(in string) string {
			in = strings.Replace(in, `
				{
					"expr": "todo_server_items_count",
					"format": "time_series",
					"instant": false,
					"intervalFactor": 1,
					"legendFormat": "items",
					"refId": "A"
				},`, ``, 1)

			return in
		},
		"cmd/server/v1/wire.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(
				in, `		items.Providers,`, "{{ range $i, $dm := .DataTypes }} {{ lower $dm.Name }}s.Providers, {{ end }}", 1)

			return in
		},
		"database/v1/database.go": func(in string) string {
			in = strings.Replace(in, `
		models.ItemDataManager`, "{{ range $i, $dm := .DataTypes }}\n\tmodels.{{ camelcase $dm.Name }}DataManager\n{{ end }}", 1)

			return in
		},
		"database/v1/queriers/postgres/migrations.go": func(in string) string {
			sentinelVal := "REPLACEMEHERE"
			replacement := `{{ range $i, $dt := .DataTypes }} 
		{
			Version:     {{ add $i 4 }} ,
			Description: "create {{ lower $dt.Name }}s table",
			Script: ` + "`" + ` 
			CREATE TABLE IF NOT EXISTS {{ lower $dt.Name }}s (
				"id" bigserial NOT NULL PRIMARY KEY,
				{{ range $j, $field := $dt.Fields }}
					"{{ snakecase $field.Name }}" {{ typeToPostgresType $field.Type }} {{ if ne true $field.Pointer }} NOT NULL {{ end }} , 
				{{ end }}				
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				"belongs_to" bigint NOT NULL,
				FOREIGN KEY ("belongs_to") REFERENCES "users"("id")
			);` + "`," + `
		},
			{{ end }}`

			in = strings.Replace(in, `		{
			Version:     4,
			Description: "create items table",
			Script: `, sentinelVal, 1)

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

			in = strings.Replace(in, sentinelVal, replacement, 1)

			return in
		},

		"database/v1/queriers/postgres/items.go": func(in string) string {
			in = strings.Replace(in, `
	itemsTableColumns = []string{
		"id",
		"name",
		"details",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}`, `
	{{ camelCase .Name }}sTableColumns = []string{
		"id",
		{{ range $j, $field := .Fields }}
			"{{ snakecase $field.Name }}", 
		{{ end }}	
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}`, 1)

			in = strings.Replace(in, `
	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}`, `
	if err := scan.Scan(
		&x.ID,
		{{ range $j, $field := .Fields }}
			&x.{{ pascal $field.Name }}, 
		{{ end }}	
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}`, 1)

			in = strings.Replace(in, `
// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (p *Postgres) buildCreateItemQuery(input *models.Item) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			"name",
			"details",
			"belongs_to",
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	logQueryBuildingError(p.logger, err)

	return query, args
}`, `
// buildCreate{{ pascal .Name }}Query takes an item and returns a creation query for that item and the relevant arguments.
func (p *Postgres) buildCreate{{ pascal .Name }}Query(input *models.{{ pascal .Name }}) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			{{ range $j, $field := .Fields }}
				{{ if $field.ValidForCreationInput }}"{{ snakecase $field.Name }}",{{ end }}
			{{ end }}
			"belongs_to",
		).
		Values(
		{{ range $j, $field := .Fields }}
			{{ if $field.ValidForCreationInput }}input.{{ pascal $field.Name }},{{ end }}
		{{ end }}	
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	logQueryBuildingError(p.logger, err)

	return query, args
}`, 1)

			in = strings.Replace(in, `
// CreateItem creates an item in the database
func (p *Postgres) CreateItem(ctx context.Context, input *models.ItemCreationInput) (*models.Item, error) {
	i := &models.Item{
		Name:      input.Name,
		Details:   input.Details,
		BelongsTo: input.BelongsTo,
	}

	query, args := p.buildCreateItemQuery(i)

	// create the item
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&i.ID, &i.CreatedOn)
	if err != nil {
		return nil, errors.Wrap(err, "error executing item creation query")
	}

	return i, nil
}`, `
// Create{{ pascal .Name }} creates an {{ lower .Name }} in the database
func (p *Postgres) Create{{ pascal .Name }}(ctx context.Context, input *models.{{ pascal .Name }}CreationInput) (*models.{{ pascal .Name }}, error) {
	x := &models.{{ pascal .Name }}{
		{{ range $j, $field := .Fields }}
			{{ if $field.ValidForCreationInput }}
				{{ pascal $field.Name }}: input.{{ pascal $field.Name }}, 
			{{ end }}
		{{ end }}
		BelongsTo: input.BelongsTo,
	}

	query, args := p.buildCreate{{ pascal .Name }}Query(x)

	// create the {{ lower .Name }}
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, errors.Wrap(err, "error executing {{ lower .Name }} creation query")
	}

	return x, nil
}`, 1)

			in = strings.Replace(in, `
		Set("name", input.Name).
		Set("details", input.Details).
`, `
		{{ range $j, $field := .Fields }}	
			{{ if $field.ValidForUpdateInput }}
				Set("{{ snakecase $field.Name }}", input.{{ pascal $field.Name }}). 
			{{ end }}	
		{{ end }}	
`, 1)

			return in
		},
		"database/v1/queriers/postgres/items_test.go": func(in string) string {
			in = strings.Replace(in, `Name: "name",`, `
		{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}	
			`, 1)

			in = strings.Replace(in, `

func buildMockRowFromItem(item *models.Item) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(itemsTableColumns).
		AddRow(
			item.ID,
			item.Name,
			item.Details,
			item.CreatedOn,
			item.UpdatedOn,
			item.ArchivedOn,
			item.BelongsTo,
		)

	return exampleRows
}

func buildErroneousMockRowFromItem(item *models.Item) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(itemsTableColumns).
		AddRow(
			item.ArchivedOn,
			item.Name,
			item.Details,
			item.CreatedOn,
			item.UpdatedOn,
			item.BelongsTo,
			item.ID,
		)

	return exampleRows
}
`, `
func buildMockRowFrom{{ pascal .Name }}({{ camelCase .Name}} *models.{{ pascal .Name }}) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows({{ camelCase .Name}}sTableColumns).
		AddRow(
			{{ camelCase .Name}}.ID,
		{{ $og := . }}{{ range $j, $field := .Fields }}	
			{{ camelCase $og.Name}}.{{ pascal $field.Name }},
		{{ end }}	
			{{ camelCase .Name}}.CreatedOn,
			{{ camelCase .Name}}.UpdatedOn,
			{{ camelCase .Name}}.ArchivedOn,
			{{ camelCase .Name}}.BelongsTo,
		)

	return exampleRows
}

func buildErroneousMockRowFrom{{ pascal .Name }}({{ camelCase .Name}} *models.{{ pascal .Name }}) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows({{ camelCase .Name}}sTableColumns).
		AddRow(
			{{ camelCase .Name}}.ArchivedOn,
		{{ $og := . }}{{ range $j, $field := .Fields }}	
			{{ camelCase $og.Name}}.{{ pascal $field.Name }},
		{{ end }}	
			{{ camelCase .Name}}.CreatedOn,
			{{ camelCase .Name}}.UpdatedOn,
			{{ camelCase .Name}}.BelongsTo,
			{{ camelCase .Name}}.ID,
		)

	return exampleRows
}
`, 1)
			in = strings.NewReplacer(
				`Name:    "name",
			Details: "details",`, `
		{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}`, `Name:      "name",
			Details:   "details",`, `
		{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}`).Replace(in)

			in = strings.Replace(in, `
func TestPostgres_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.Item{
			Name:      "name",
			Details:   "details",
			BelongsTo: 123,
		}

		expectedArgCount := 3
		expectedQuery := "INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"

		actualQuery, args := p.buildCreateItemQuery(expected)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Details, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
	})
}
`, `
func TestPostgres_buildCreate{{ pascal .Name }}Query(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.{{ pascal .Name }}{
		{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}	
			BelongsTo: 123,
		}

		// NOTE: this test is a deliberate failure

		expectedArgCount := 0
		expectedQuery := "INSERT INTO {{ snakecase .Name }} () VALUES () RETURNING id, created_on"

		actualQuery, args := p.buildCreate{{ pascal .Name }}Query(expected)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)

		{{ range $j, $field := .Fields }}	
		assert.Equal(t, expected.{{ pascal $field.Name }}, args[{{ $j }}].({{ $field.Type }}))
		{{ end }}	
		assert.Equal(t, expected.BelongsTo, args[len(args)-1].(uint64))
	})
}`, 1)

			in = strings.Replace(in, `
func TestPostgres_CreateItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ItemCreationInput{
			Name:      expected.Name,
			BelongsTo: expected.BelongsTo,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).
			AddRow(expected.ID, uint64(time.Now().Unix()))

		expectedQuery := "INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Name,
				expected.Details,
				expected.BelongsTo,
			).
			WillReturnRows(exampleRows)

		actual, err := p.CreateItem(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ItemCreationInput{
			Name:      example.Name,
			BelongsTo: example.BelongsTo,
		}

		expectedQuery := "INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.Name,
				example.Details,
				example.BelongsTo,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.CreateItem(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`, `
func TestPostgres_Create{{ pascal .Name }}(T *testing.T) {
	T.Parallel()
	
	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.{{ pascal .Name }}{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.{{ pascal .Name }}CreationInput{
			{{ range $j, $field := .Fields }}
				{{ if $field.ValidForCreationInput }}
					{{ pascal $field.Name }}: expected.{{ pascal $field.Name }},
				{{ end }}
			{{ end }}	
			BelongsTo: expected.BelongsTo,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).
			AddRow(expected.ID, uint64(time.Now().Unix()))

		expectedQuery := "INSERT INTO {{ snakecase .Name }} (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
			{{ $og := . }}{{ range $j, $field := .Fields }}	
				expected.{{ pascal $field.Name }},
			{{ end }}	
				expected.BelongsTo,
			).
			WillReturnRows(exampleRows)

		actual, err := p.Create{{ pascal .Name }}(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.{{ pascal .Name }}{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.{{ pascal .Name }}CreationInput{
			BelongsTo: example.BelongsTo,
		}

		expectedQuery := "INSERT INTO {{ snakecase .Name }} (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
			{{ $og := . }}{{ range $j, $field := .Fields }}	
				example.{{ pascal $field.Name }},
			{{ end }}	
				example.BelongsTo,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.Create{{ pascal .Name }} (context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}`, 1)

			in = strings.ReplaceAll(in, `name, details, `, "")
			in = strings.ReplaceAll(in, `name,details,`, "")
			in = strings.ReplaceAll(in, `name = $1, details = $2, `, "")

			in = strings.ReplaceAll(in, `
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),`, `
			ID:        123,
			{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
			{{ end }}	
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),`)

			in = strings.ReplaceAll(in, `WithArgs(
				example.Name,
				example.Details,
				example.BelongsTo,
				example.ID,
			).`, `WithArgs(
			{{ range $j, $field := .Fields }}	
				example.{{ pascal $field.Name }},
			{{ end }}	
				example.BelongsTo,
				example.ID,
			).`)

			in = strings.ReplaceAll(in, `
			WithArgs(
				expected.Name,
				expected.Details,
				expected.BelongsTo,
			).`, `
			WithArgs(
			{{ range $j, $field := .Fields }}	
				expected.{{ pascal $field.Name }},
			{{ end }}	
				expected.BelongsTo,
			).`)

			in = strings.ReplaceAll(in, `
			WithArgs(
				example.Name,
				example.Details,
				example.BelongsTo,
			).`, `
			WithArgs(
			{{ range $j, $field := .Fields }}	
				example.{{ pascal $field.Name }},
			{{ end }}	
				example.BelongsTo,
			).`)

			in = strings.ReplaceAll(in, `
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Details, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))`, `
		{{ range $j, $field := .Fields }}	
		assert.Equal(t, expected.{{ pascal $field.Name }}, args[{{ $j }}].({{ $field.Type }}))
		{{ end }}	
		assert.Equal(t, expected.BelongsTo, args[len(args)-1].(uint64))`)

			in = strings.Replace(in, `
			WithArgs(
				expected.Name,
				expected.Details,
				expected.BelongsTo,
				expected.ID,
			).`, `
			WithArgs(
				{{ range $j, $field := .Fields }}	
					expected.{{ pascal $field.Name }},
				{{ end }}	
				expected.BelongsTo,
				expected.ID,
			).`, 1)

			in = strings.ReplaceAll(in, `&models.Item{
			Name: "name",
		}`, `&models.Item{
			{{ range $j, $field := .Fields }}	
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
			{{ end }}	
		}`)

			return in
		},

		"models/v1/item.go": func(in string) string {
			z := `
	// Item represents an item
	Item struct {
		ID         uint64  ` + "`" + `json:"id"` + "`" + `
		Name       string  ` + "`" + `json:"name"` + "`" + `
		Details    string  ` + "`" + `json:"details"` + "`" + `
		CreatedOn  uint64  ` + "`" + `json:"created_on"` + "`" + `
		UpdatedOn  *uint64 ` + "`" + `json:"updated_on"` + "`" + `
		ArchivedOn *uint64 ` + "`" + `json:"archived_on"` + "`" + `
		BelongsTo  uint64  ` + "`" + `json:"belongs_to"` + "`" + `
	}`
			u := `
	// {{ camelcase .Name }} represents a(n) {{ camelcase .Name }}
	{{ camelcase .Name }} struct {
		ID         uint64  ` + "`" + `json:"id"` + "`" + `
		{{ range $i, $field := .Fields }}
			{{ pascal $field.Name }} {{ if $field.Pointer }}*{{ end }}{{ $field.Type }} ` + "`" + `json:"{{ snakecase $field.Name }}"` + "`" + `
		{{ end }}
		CreatedOn  uint64  ` + "`" + `json:"created_on"` + "`" + `
		UpdatedOn  *uint64 ` + "`" + `json:"updated_on"` + "`" + `
		ArchivedOn *uint64 ` + "`" + `json:"archived_on"` + "`" + `
		BelongsTo  uint64  ` + "`" + `json:"belongs_to"` + "`" + `
	}`
			in = strings.ReplaceAll(in, z, u)

			in = strings.Replace(in, `
	// ItemCreationInput represents what a user could set as input for creating items
	ItemCreationInput struct {
		Name      string `+"`"+`json:"name"`+"`"+`
		Details   string `+"`"+`json:"details"`+"`"+`
		BelongsTo uint64 `+"`"+`json:"-"`+"`"+`
	}

	// ItemUpdateInput represents what a user could set as input for updating items
	ItemUpdateInput struct {
		Name      string `+"`"+`json:"name"`+"`"+`
		Details   string `+"`"+`json:"details"`+"`"+`
		BelongsTo uint64 `+"`"+`json:"-"`+"`"+`
	}`, `
	// {{ camelcase .Name }}CreationInput represents what a user could set as input for creating {{ camelcase .Name }}s
	{{ camelcase .Name }}CreationInput struct {
		{{ range $i, $field := .Fields }}
			{{ if $field.ValidForCreationInput }}
				{{ $field.Name }} {{ if $field.Pointer }}*{{ end }}{{ $field.Type }} `+"`"+`json:"{{ snakecase $field.Name }}"`+"`"+`
			{{ end }}
		{{ end }}
		BelongsTo uint64 `+"`"+`json:"-"`+"`"+`
	}

	// {{ camelcase .Name }}UpdateInput represents what a user could set as input for updating {{ camelcase .Name }}s
	{{ camelcase .Name }}UpdateInput struct {
		{{ range $i, $field := .Fields }}
			{{ if $field.ValidForUpdateInput }}
				{{ $field.Name }} {{ if $field.Pointer }}*{{ end }}{{ $field.Type }} `+"`"+`json:"{{ snakecase $field.Name }}"`+"`"+`
			{{ end }}
		{{ end }}
		BelongsTo uint64 `+"`"+`json:"-"`+"`"+`
	}`, 1)

			in = strings.Replace(in, `if input.Name != "" || input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Details != "" || input.Details != x.Details {
		x.Details = input.Details
	}`, `{{ range $i, $field := .Fields }}
		{{ if $field.ValidForUpdateInput }} 
		if x.{{ $field.Name }} != input.{{ $field.Name }} {
			x.{{ $field.Name }} = input.{{ $field.Name }} 
		}{{ end }}
{{ end }}`, 1)

			return in
		},
		"server/v1/http/server.go": func(in string) string {
			in = strings.Replace(in, "		itemsService         models.ItemDataServer", `
{{ range $i, $dt := .DataTypes }} 
	{{ lower $dt.Name }}sService models.{{ camelcase $dt.Name }}DataServer
{{ end }}`, 1)
			in = strings.Replace(in, "	itemsService models.ItemDataServer,", `
{{ range $i, $dt := .DataTypes }} 
				{{ lower $dt.Name }}sService models.{{ camelcase $dt.Name }}DataServer, 
{{ end }}`, 1)
			in = strings.Replace(in, "		itemsService:         itemsService,", `
{{ range $i, $dt := .DataTypes }} 
	{{ lower $dt.Name }}sService: {{ lower $dt.Name }}sService,
{{ end }}`, 1)

			return in
		},
		"server/v1/http/server_test.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)
			in = strings.Replace(in, "		itemsService:         &mmodels.ItemDataServer{},", `
{{ range $i, $dt := .DataTypes }} 
	{{ lower $dt.Name }}sService: &mmodels.{{ camelcase $dt.Name }}DataServer{}, 
{{ end }}`, 1)
			in = strings.Replace(in, "			&items.Service{},", `
{{ range $i, $dt := .DataTypes }} 
	&{{ lower $dt.Name }}s.Service{}, 
{{ end }}`, 1)

			return in
		},
		"server/v1/http/wire_param_fetchers.go": func(in string) string {
			in = strings.Replace(
				in, `	"{{ .OutputRepository }}/services/v1/items"`,
				iterablesImportsTemplateCode, 1,
			)

			in = strings.Replace(in, "		ProvideItemIDFetcher,", `
{{ range $i, $dt := .DataTypes }}
	Provide{{ camelcase $dt.Name }}IDFetcher,
{{ end }}
`, 1)

			in = strings.Replace(in, `// ProvideItemIDFetcher provides an ItemIDFetcher
func ProvideItemIDFetcher(logger logging.Logger) items.ItemIDFetcher {
	return buildChiItemIDFetcher(logger)
}`, `
{{ range $i, $dt := .DataTypes }}
// Provide{{ camelcase $dt.Name }}IDFetcher provides an {{ camelcase $dt.Name }}IDFetcher
func Provide{{ camelcase $dt.Name }}IDFetcher(logger logging.Logger) {{ lower $dt.Name }}s.{{ camelcase $dt.Name }}IDFetcher {
	return buildChi{{ camelcase $dt.Name }}IDFetcher(logger)
}
{{ end }}`, 1)

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
}`, `
{{ range $i, $dt := .DataTypes }}
// chi{{ camelcase $dt.Name }}IDFetcher fetches a {{ camelcase $dt.Name }}ID from a request routed by chi.
func buildChi{{ camelcase $dt.Name }}IDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, {{ lower $dt.Name }}s.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching {{ camelcase $dt.Name }}ID from request")
		}
		return u
	}
}
{{ end }}`, 1)

			in = strings.Replace(in, `
// ProvideUserIDFetcher provides a UserIDFetcher
func ProvideUserIDFetcher() items.UserIDFetcher {
	return UserIDFetcher
}`, `
{{ range $i, $dt := .DataTypes }}
// Provide{{ camelcase $dt.Name }}UserIDFetcher provides a UserIDFetcher
func Provide{{ camelcase $dt.Name }}UserIDFetcher() {{ lower $dt.Name }}s.UserIDFetcher {
	return UserIDFetcher
}
{{ end }}`, 1)

			in = strings.Replace(in, `		ProvideUserIDFetcher,`, `
				{{ range $i, $dt := .DataTypes }}
					Provide{{ camelcase $dt.Name }}UserIDFetcher, 
				{{ end }}
			`, 1)

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
`, `{{ range $i, $dt := .DataTypes }}
func TestProvide{{ camelcase $dt.Name }}UserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = Provide{{ camelcase $dt.Name }}UserIDFetcher()
	})
}
{{ end }}`, 1)

			in = strings.Replace(in, `func TestProvideItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemIDFetcher(noop.ProvideNoopLogger())
	})
}`, `{{ range $i, $dt := .DataTypes }} 
func TestProvide{{ camelcase $dt.Name }}IDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = Provide{{ camelcase $dt.Name }}IDFetcher(noop.ProvideNoopLogger())
	})
}
{{ end }}`, 1)

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
}`,
				`
{{ range $i, $dt := .DataTypes }} 
func Test_buildChi{{ camelcase $dt.Name }}IDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildChi{{ camelcase $dt.Name }}IDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{ {{ lower $dt.Name }}s.URIParamKey },
				Values: []string{fmt.Sprintf("%d", expected)},
			},
		}))

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildChi{{ camelcase $dt.Name }}IDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{ {{ lower $dt.Name }}s.URIParamKey },
				Values: []string{"expected"},
			},
		}))

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
{{ end }}`,
				1)

			return in
		},
		"services/v1/frontend/http_routes.go": func(in string) string {
			in = strings.Replace(in, `
	// itemsFrontendPathRegex matches URLs against our frontend router's specification for specific item routes
	itemsFrontendPathRegex = regexp.MustCompile(`+"`"+`/items/\d+`+"`)", `
{{ range $i, $dt := .DataTypes }}
	// {{ lower $dt.Name }}sFrontendPathRegex matches URLs against our frontend router's specification for specific {{ camelcase $dt.Name }} routes
	{{ lower $dt.Name }}sFrontendPathRegex = regexp.MustCompile(`+"`"+`/items/\d+`+"`"+`)
{{ end }}
`, 1)
			in = strings.Replace(in, `

			"/items/new",`, "", 1)
			in = strings.Replace(in, `
		if itemsFrontendPathRegex.MatchString(req.URL.Path) {
			req.URL.Path = "/"
		}`, `{{ range $i, $dt := .DataTypes }}
	if {{ lower $dt.Name }}sFrontendPathRegex.MatchString(req.URL.Path) {
		req.URL.Path = "/"
	}
{{ end }}`, 1)

			return in
		},
		"services/v1/frontend/http_routes_test.go": func(in string) string {
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
`, `
{{ range $i, $dt := .DataTypes }}
	T.Run("with frontend {{ snakecase $dt.Name }}s routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/{{ snakecase $dt.Name }}s/123"

		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
{{ end }}`, 1)

			in = strings.ReplaceAll(in, `Name:    "name",`, `
`)

			in = strings.Replace(in, `
			Name:    expected.Name,
			Details: expected.Details,
`, `{{ range $i, $dt := .DataTypes }}
	{{ range $j, $field := $dt.Fields }}
		{{ $field.Name }}: expected.{{ $field.Name}},
	{{ end }}
{{ end }}`, -1)

			return in
		},
		"database/v1/database_mock.go": func(in string) string {
			in = strings.Replace(in, `
		ItemDataManager:         &mmodels.ItemDataManager{},`, "\n{{ range $i, $dt := .DataTypes }}\n {{ camelcase $dt.Name }}DataManager: &mmodels.{{ camelcase $dt.Name }}DataManager{},\n{{ end }}", 1)
			in = strings.Replace(in, `
	*mmodels.ItemDataManager`, "\n{{ range $i, $dt := .DataTypes }}\n *mmodels.{{ camelcase $dt.Name }}DataManager \n {{ end }}", 1)

			return in
		},
		"frontend/v1/src/App.svelte": func(in string) string {
			in = strings.Replace(in, `

  // Items routes
  import ReadItem from "./pages/items/Read.svelte";
  import CreateItem from "./pages/items/Create.svelte";
  import Items from "./pages/items/List.svelte";`, `{{ range $i, $dt := .DataTypes }}
  // {{ camelcase $dt.Name }} routes
  import Read{{ camelcase $dt.Name }} from "./pages/{{ lower $dt.Name }}s/Read.svelte";
  import Create{{ camelcase $dt.Name }} from "./pages/{{ lower $dt.Name }}s/Create.svelte";
  import {{ camelcase $dt.Name }}s from "./pages/{{ lower $dt.Name }}s/List.svelte";
{{ end }}`, 1)

			in = strings.Replace(in, `
    <Link to="items">Items</Link>
    <Link to="items/new">Create Item</Link>`, `{{ range $i, $dt := .DataTypes }}
	<Link to="{{ lower $dt.Name }}s">{{ camelcase $dt.Name }}</Link>
	<Link to="{{ lower $dt.Name }}s/new">Create {{ camelcase $dt.Name }}</Link>
{{ end }}`, 1)

			in = strings.Replace(in, `
    <Route path="items" component={Items} />
    <Route path="items/:id" component={ReadItem} />
    <Route path="items/new" component={CreateItem} />`, `{{ range $i, $dt := .DataTypes }}
    <Route path="{{ lower $dt.Name }}s" component={ {{ camelcase $dt.Name }}s } />
    <Route path="{{ lower $dt.Name }}s/:id" component={ Read{{ camelcase $dt.Name }} } />
    <Route path="{{ lower $dt.Name }}s/new" component={ Create{{ camelcase $dt.Name }} } />
{{ end }}`, 1)

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
		"services/v1/users/http_routes.go": func(in string) string {
			in = strings.Replace(in, `
		// "otpauth://totp/{{ .Issuer }}:{{ .Username }}?secret={{ .Secret }}&issuer={{ .Issuer }}"`, "", 1)

			return in
		},
		"tests/v1/load/actions.go": func(in string) string {
			in = strings.Replace(in, `

	for k, v := range buildItemActions(c) {
		allActions[k] = v
	}`, `{{ range $i, $dt := .DataTypes }}
	for k, v := range build{{ camelcase $dt.Name }}Actions(c) {
		allActions[k] = v
	}
{{ end }}`, 1)
			return in
		},
		"tests/v1/integration/items_test.go": func(in string) string {
			in = strings.ReplaceAll(in, `"github.com/icrowley/fake"`, "")

			in = strings.ReplaceAll(in, `&models.ItemCreationInput{
					Name:    expected.Name,
					Details: expected.Details,
				})`, `&models.{{ pascal .Name }}CreationInput{
					{{ range $i, $field := .Fields }}
					{{ pascal $field.Name }}: expected.{{ pascal $field.Name }},
					{{ end }}
				})`)
			in = strings.ReplaceAll(in, `&models.Item{
					Name:    expected.Name,
					Details: expected.Details,
				})`, `&models.{{ pascal .Name }}{
					{{ range $i, $field := .Fields }}
					{{ pascal $field.Name }}: expected.{{ pascal $field.Name }},
					{{ end }}
				})`)
			in = strings.ReplaceAll(in, `&models.Item{
				Name:    "name",
				Details: "details",
			}`, `&models.{{ pascal .Name }}{
					{{ range $i, $field := .Fields }}
					{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
					{{ end }}
				}`)

			in = strings.ReplaceAll(in, `&models.Item{
				Name:    "new name",
				Details: "new details",
			}`, `&models.{{ pascal .Name }}{
			{{ range $i, $field := .Fields }}
			{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
			{{ end }}
		}`)

			in = strings.ReplaceAll(in, `&models.ItemCreationInput{
					Name:    "old name",
					Details: "old details",
				},`, `&models.{{ pascal .Name }}CreationInput{
					{{ range $i, $field := .Fields }}
					{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
					{{ end }}
				},`)

			in = strings.ReplaceAll(in, `
func checkItemEquality(t *testing.T, expected, actual *models.Item) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Details, actual.Details)
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummyItem(t *testing.T) *models.Item {
	t.Helper()

	x := &models.ItemCreationInput{
		Name:    fake.Word(),
		Details: fake.Sentence(),
	}
	y, err := todoClient.CreateItem(context.Background(), x)
	require.NoError(t, err)
	return y
}`, `
func check{{ pascal .Name }}Equality(t *testing.T, expected, actual *models.{{ pascal .Name }}) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	{{ range $i, $field := .Fields }}
	assert.Equal(t, expected.{{ pascal $field.Name }}, actual.{{ pascal $field.Name }})
	{{ end }}
	assert.NotZero(t, actual.CreatedOn)
}

func buildDummy{{ pascal .Name }}(t *testing.T) *models.{{ pascal .Name }} {
	t.Helper()

	x := &models.{{ pascal .Name }}CreationInput{
		{{ range $i, $field := .Fields }}
		{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}
	}
	y, err := todoClient.Create{{ pascal .Name }}(context.Background(), x)
	require.NoError(t, err)
	return y
}
`)

			in = strings.ReplaceAll(in, `premade.Name, premade.Details = expected.Name, expected.Details`, "// CHANGEME")

			return in
		},
		"tests/v1/testutil/rand/model/items.go": func(in string) string {
			in = strings.ReplaceAll(in, `	"github.com/icrowley/fake"`, "")

			in = strings.ReplaceAll(in, `
		Name:    fake.Word(),
		Details: fake.Sentence(),`, `
		{{ range $i, $field := .Fields }}
		{{ pascal $field.Name }}: {{ typeExample $field.Type $field.Pointer }},
		{{ end }}`)

			return in
		},
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

			if _, ok := iterableFiles[relativePath]; ok {
				outputFilepath = strings.Replace(outputFilepath, "base_repository", "iterables", 1)
				outputFilepath = strings.Replace(outputFilepath, "item", "model", 1)
				if err := os.MkdirAll(filepath.Dir(outputFilepath), os.ModePerm); err != nil {
					return err
				}

				file = strings.ReplaceAll(file, "Items", "{{ camelcase .Name }}s")
				file = strings.ReplaceAll(file, "Item", "{{ camelcase .Name }}")
				file = strings.ReplaceAll(file, "items", "{{ lower .Name }}s")
				file = strings.ReplaceAll(file, "item", "{{ lower .Name }}")
			}

			if !strings.HasSuffix(outputFilepath, "dashboard.json") {
				outputFilepath += ".tmpl"
			}
			e := ioutil.WriteFile(outputFilepath, []byte(file), info.Mode())
			return e
		}
	}); err != nil {
		log.Fatal(err)
	}
}
