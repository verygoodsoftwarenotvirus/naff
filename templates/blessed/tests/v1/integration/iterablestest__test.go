package integration

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func Test_buildRequisiteCreationCode(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildRequisiteCreationCode(proj, cherry)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
