package queriers

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/stretchr/testify/assert"
)

func Test_buildGetSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:          wordsmith.FromSingularPascalCase("Something"),
			BelongsToUser: true,
		}

		f := jen.NewFile(dbv.SingularPackageName())
		lines := buildGetSomethingQueryFuncDecl(dbv, dt)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package postgres

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id":         somethingID,
			"belongs_to": userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`

		assert.Equal(t, expected, actual)
	})

	T.Run("belonging to anther data type", func(t *testing.T) {
		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Something"),
			BelongsToUser:   false,
			BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
		}

		f := jen.NewFile(dbv.SingularPackageName())
		lines := buildGetSomethingQueryFuncDecl(dbv, dt)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package postgres

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID, ownerID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id":               somethingID,
			"belongs_to_owner": ownerID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`

		assert.Equal(t, expected, actual)
	})

	T.Run("belonging to nothing", func(t *testing.T) {
		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Something"),
			BelongsToUser:   false,
			BelongsToStruct: nil,
			BelongsToNobody: true,
		}

		f := jen.NewFile(dbv.SingularPackageName())
		lines := buildGetSomethingQueryFuncDecl(dbv, dt)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package postgres

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id": somethingID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTableColumns(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		dt := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Something"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("FieldOne"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldTwo"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldThree"),
				},
			},
			BelongsToUser: true,
		}

		actual := buildTableColumns(dt)
		expected := []jen.Code{
			jen.Lit("id"),
			jen.Lit("field_one"),
			jen.Lit("field_two"),
			jen.Lit("field_three"),
			jen.Lit("created_on"),
			jen.Lit("updated_on"),
			jen.Lit("archived_on"),
			jen.Lit("belongs_to_user"),
		}

		assert.Equal(t, actual, expected)
	})

	T.Run("with alternative ownership", func(t *testing.T) {
		dt := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Something"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("FieldOne"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldTwo"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldThree"),
				},
			},
			BelongsToStruct: wordsmith.FromSingularPascalCase("Another"),
		}

		actual := buildTableColumns(dt)
		expected := []jen.Code{
			jen.Lit("id"),
			jen.Lit("field_one"),
			jen.Lit("field_two"),
			jen.Lit("field_three"),
			jen.Lit("created_on"),
			jen.Lit("updated_on"),
			jen.Lit("archived_on"),
			jen.Lit("belongs_to_another"),
		}

		assert.Equal(t, actual, expected)
	})

	T.Run("with no ownership", func(t *testing.T) {
		dt := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Something"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("FieldOne"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldTwo"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("FieldThree"),
				},
			},
			BelongsToUser:   false,
			BelongsToStruct: nil,
			BelongsToNobody: true,
		}

		actual := buildTableColumns(dt)
		expected := []jen.Code{
			jen.Lit("id"),
			jen.Lit("field_one"),
			jen.Lit("field_two"),
			jen.Lit("field_three"),
			jen.Lit("created_on"),
			jen.Lit("updated_on"),
			jen.Lit("archived_on"),
		}

		assert.Equal(t, actual, expected)
	})
}
