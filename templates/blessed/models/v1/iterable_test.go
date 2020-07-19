package v1

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"
)

var (
	a = models.DataType{
		Name: wordsmith.FromSingularPascalCase("Grandparent"),
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("GrandparentName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
	b = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Parent"),
		BelongsToStruct: a.Name,
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("ParentName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
	c = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Child"),
		BelongsToStruct: b.Name,
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("ChildName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
)

func Test_buildInterfaceMethods(T *testing.T) {
	T.Parallel()

	T.Run("lone type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		code := jen.NewFile("farts")
		code.Add(
			jen.Type().ID("SomethingDataManager").Interface(
				buildInterfaceMethods(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := code.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

type SomethingDataManager interface {
	GetChild(ctx context.Context, grandparentID, parentID, childID uint64) (*Child, error)
	GetChildCount(ctx context.Context, grandparentID, parentID uint64, filter *QueryFilter) (uint64, error)
	GetAllChildrenCount(ctx context.Context) (uint64, error)
	GetChildren(ctx context.Context, grandparentID, parentID uint64, filter *QueryFilter) (*ChildList, error)
	GetAllChildrenForParent(ctx context.Context, grandparentID, parentID uint64) ([]Child, error)
	CreateChild(ctx context.Context, grandparentID, parentID uint64, input *ChildCreationInput) (*Child, error)
	UpdateChild(ctx context.Context, grandparentID uint64, updated *Child) error
	ArchiveChild(ctx context.Context, grandparentID, parentID, childID uint64) error
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
