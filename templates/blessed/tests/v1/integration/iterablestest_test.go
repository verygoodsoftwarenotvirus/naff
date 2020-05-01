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

func Test_buildBuildDummySomething(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{DataTypes: []models.DataType{a, b, c}}

		ret := jen.NewFile("farts")
		ret.Add(buildBuildDummySomething(proj, c)...)

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	require "github.com/stretchr/testify/require"
	v1 "models/v1"
	"testing"
)

func buildDummyChild(t *testing.T, ctx context.Context) *v1.Child {
	t.Helper()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	x := &v1.ChildCreationInput{
		ChildName: gofakeit.Word(),
	}
	y, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, x)
	require.NoError(t, err)

	return y
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("one dependency", func(t *testing.T) {
		proj := &models.Project{DataTypes: []models.DataType{a, b, c}}

		ret := jen.NewFile("farts")
		ret.Add(buildBuildDummySomething(proj, b)...)

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	require "github.com/stretchr/testify/require"
	v1 "models/v1"
	"testing"
)

func buildDummyParent(t *testing.T, ctx context.Context) *v1.Parent {
	t.Helper()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	x := &v1.ParentCreationInput{
		ParentName: gofakeit.Word(),
	}
	y, err := todoClient.CreateParent(ctx, createdGrandparent.ID, x)
	require.NoError(t, err)

	return y
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("lone type", func(t *testing.T) {
		proj := &models.Project{DataTypes: []models.DataType{a, b, c}}

		ret := jen.NewFile("farts")
		ret.Add(buildBuildDummySomething(proj, a)...)

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	require "github.com/stretchr/testify/require"
	v1 "models/v1"
	"testing"
)

func buildDummyGrandparent(t *testing.T, ctx context.Context) *v1.Grandparent {
	t.Helper()

	x := &v1.GrandparentCreationInput{
		GrandparentName: gofakeit.Word(),
	}
	y, err := todoClient.CreateGrandparent(ctx, x)
	require.NoError(t, err)

	return y
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildRequisiteCreationCode(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{DataTypes: []models.DataType{a, b, c}}

		ret := jen.NewFile("farts")
		ret.Add(jen.Func().ID("doSomething").Params().Block(
			buildRequisiteCreationCode(proj, c)...,
		))

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	gofakeit "github.com/brianvoe/gofakeit"
	v1 "models/v1"
)

func doSomething() {
	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create child
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	checkValueAndError(t, createdChild, err)

}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildRequisiteCleanupCode(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{DataTypes: []models.DataType{a, b, c}}

		ret := jen.NewFile("farts")
		ret.Add(jen.Func().ID("doSomething").Params().Block(
			buildRequisiteCleanupCode(proj, c)...,
		))

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	assert "github.com/stretchr/testify/assert"
)

func doSomething() {

	// Clean up child
	assert.NoError(t, todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID))

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildCreationArguments(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildCreationArguments(proj, "expected", c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import ()

func doSomething() {
	expectedGrandparent.ID
	expectedParent.ID
	expectedChild.ID
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildEqualityCheckLines(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildEqualityCheckLines(
					models.DataType{
						Name:            wordsmith.FromSingularPascalCase("Everything"),
						BelongsToUser:   false,
						BelongsToNobody: true,
						BelongsToStruct: nil,
						Fields: []models.DataField{
							{
								Name: wordsmith.FromSingularPascalCase("stringVar"),
								Type: "string",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToStringVar"),
								Type:    "string",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("boolVar"),
								Type: "bool",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToBoolVar"),
								Type:    "bool",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("intVar"),
								Type: "int",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToIntVar"),
								Type:    "int",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("int8Var"),
								Type: "int8",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToint8Var"),
								Type:    "int8",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("int16Var"),
								Type: "int16",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToInt16Var"),
								Type:    "int16",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("int32Var"),
								Type: "int32",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToInt32Var"),
								Type:    "int32",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("int64Var"),
								Type: "int64",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToInt64Var"),
								Type:    "int64",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("uintVar"),
								Type: "uint",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToUintVar"),
								Type:    "uint",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("uint8Var"),
								Type: "uint8",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToUint8Var"),
								Type:    "uint8",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("uint16Var"),
								Type: "uint16",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToUint16Var"),
								Type:    "uint16",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("uint32Var"),
								Type: "uint32",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToUint32Var"),
								Type:    "uint32",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("uint64Var"),
								Type: "uint64",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToUint64Var"),
								Type:    "uint64",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("float32Var"),
								Type: "float32",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToFloat32Var"),
								Type:    "float32",
								Pointer: true,
							},
							{
								Name: wordsmith.FromSingularPascalCase("float64Var"),
								Type: "float64",
							},
							{
								Name:    wordsmith.FromSingularPascalCase("pointerToFloat64Var"),
								Type:    "float64",
								Pointer: true,
							},
						},
					},
				)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	assert "github.com/stretchr/testify/assert"
)

func doSomething() {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.StringVar, actual.StringVar, "expected StringVar for ID %d to be %v, but it was %v ", expected.ID, expected.StringVar, actual.StringVar)
	assert.Equal(t, *expected.PointerToStringVar, *actual.PointerToStringVar, "expected PointerToStringVar to be %v, but it was %v ", expected.PointerToStringVar, actual.PointerToStringVar)
	assert.Equal(t, expected.BoolVar, actual.BoolVar, "expected BoolVar for ID %d to be %v, but it was %v ", expected.ID, expected.BoolVar, actual.BoolVar)
	assert.Equal(t, *expected.PointerToBoolVar, *actual.PointerToBoolVar, "expected PointerToBoolVar to be %v, but it was %v ", expected.PointerToBoolVar, actual.PointerToBoolVar)
	assert.Equal(t, expected.IntVar, actual.IntVar, "expected IntVar for ID %d to be %v, but it was %v ", expected.ID, expected.IntVar, actual.IntVar)
	assert.Equal(t, *expected.PointerToIntVar, *actual.PointerToIntVar, "expected PointerToIntVar to be %v, but it was %v ", expected.PointerToIntVar, actual.PointerToIntVar)
	assert.Equal(t, expected.Int8Var, actual.Int8Var, "expected Int8Var for ID %d to be %v, but it was %v ", expected.ID, expected.Int8Var, actual.Int8Var)
	assert.Equal(t, *expected.PointerToint8Var, *actual.PointerToint8Var, "expected PointerToint8Var to be %v, but it was %v ", expected.PointerToint8Var, actual.PointerToint8Var)
	assert.Equal(t, expected.Int16Var, actual.Int16Var, "expected Int16Var for ID %d to be %v, but it was %v ", expected.ID, expected.Int16Var, actual.Int16Var)
	assert.Equal(t, *expected.PointerToInt16Var, *actual.PointerToInt16Var, "expected PointerToInt16Var to be %v, but it was %v ", expected.PointerToInt16Var, actual.PointerToInt16Var)
	assert.Equal(t, expected.Int32Var, actual.Int32Var, "expected Int32Var for ID %d to be %v, but it was %v ", expected.ID, expected.Int32Var, actual.Int32Var)
	assert.Equal(t, *expected.PointerToInt32Var, *actual.PointerToInt32Var, "expected PointerToInt32Var to be %v, but it was %v ", expected.PointerToInt32Var, actual.PointerToInt32Var)
	assert.Equal(t, expected.Int64Var, actual.Int64Var, "expected Int64Var for ID %d to be %v, but it was %v ", expected.ID, expected.Int64Var, actual.Int64Var)
	assert.Equal(t, *expected.PointerToInt64Var, *actual.PointerToInt64Var, "expected PointerToInt64Var to be %v, but it was %v ", expected.PointerToInt64Var, actual.PointerToInt64Var)
	assert.Equal(t, expected.UintVar, actual.UintVar, "expected UintVar for ID %d to be %v, but it was %v ", expected.ID, expected.UintVar, actual.UintVar)
	assert.Equal(t, *expected.PointerToUintVar, *actual.PointerToUintVar, "expected PointerToUintVar to be %v, but it was %v ", expected.PointerToUintVar, actual.PointerToUintVar)
	assert.Equal(t, expected.Uint8Var, actual.Uint8Var, "expected Uint8Var for ID %d to be %v, but it was %v ", expected.ID, expected.Uint8Var, actual.Uint8Var)
	assert.Equal(t, *expected.PointerToUint8Var, *actual.PointerToUint8Var, "expected PointerToUint8Var to be %v, but it was %v ", expected.PointerToUint8Var, actual.PointerToUint8Var)
	assert.Equal(t, expected.Uint16Var, actual.Uint16Var, "expected Uint16Var for ID %d to be %v, but it was %v ", expected.ID, expected.Uint16Var, actual.Uint16Var)
	assert.Equal(t, *expected.PointerToUint16Var, *actual.PointerToUint16Var, "expected PointerToUint16Var to be %v, but it was %v ", expected.PointerToUint16Var, actual.PointerToUint16Var)
	assert.Equal(t, expected.Uint32Var, actual.Uint32Var, "expected Uint32Var for ID %d to be %v, but it was %v ", expected.ID, expected.Uint32Var, actual.Uint32Var)
	assert.Equal(t, *expected.PointerToUint32Var, *actual.PointerToUint32Var, "expected PointerToUint32Var to be %v, but it was %v ", expected.PointerToUint32Var, actual.PointerToUint32Var)
	assert.Equal(t, expected.Uint64Var, actual.Uint64Var, "expected Uint64Var for ID %d to be %v, but it was %v ", expected.ID, expected.Uint64Var, actual.Uint64Var)
	assert.Equal(t, *expected.PointerToUint64Var, *actual.PointerToUint64Var, "expected PointerToUint64Var to be %v, but it was %v ", expected.PointerToUint64Var, actual.PointerToUint64Var)
	assert.Equal(t, expected.Float32Var, actual.Float32Var, "expected Float32Var for ID %d to be %v, but it was %v ", expected.ID, expected.Float32Var, actual.Float32Var)
	assert.Equal(t, *expected.PointerToFloat32Var, *actual.PointerToFloat32Var, "expected PointerToFloat32Var to be %v, but it was %v ", expected.PointerToFloat32Var, actual.PointerToFloat32Var)
	assert.Equal(t, expected.Float64Var, actual.Float64Var, "expected Float64Var for ID %d to be %v, but it was %v ", expected.ID, expected.Float64Var, actual.Float64Var)
	assert.Equal(t, *expected.PointerToFloat64Var, *actual.PointerToFloat64Var, "expected PointerToFloat64Var to be %v, but it was %v ", expected.PointerToFloat64Var, actual.PointerToFloat64Var)
	assert.NotZero(t, actual.CreatedOn)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestCreating(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestCreating(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create child
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	checkValueAndError(t, createdChild, err)

	// Assert child equality
	checkChildEquality(t, createdChild, exampleChild)

	// Clean up
	err = todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID)
	assert.NoError(t, err)

	actual, err := todoClient.GetChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID)
	checkValueAndError(t, actual, err)
	checkChildEquality(t, createdChild, actual)
	assert.NotZero(t, actual.ArchivedOn)

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestListing(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestListing(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create children
	var expected []*v1.Child
	for i := 0; i < 5; i++ {
		exampleChild := &v1.Child{
			ChildName: gofakeit.Word(),
		}

		createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
			ChildName: exampleChild.ChildName,
		})
		checkValueAndError(t, createdChild, err)

		expected = append(expected, createdChild)
	}

	// Assert child list equality
	actual, err := todoClient.GetChildren(ctx, createdGrandparent.ID, createdParent.ID, nil)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual.Children),
		"expected %d to be <= %d",
		len(expected),
		len(actual.Children),
	)

	// Clean up
	for _, createdChild := range actual.Children {
		err = todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID)
		assert.NoError(t, err)
	}

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("single dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestListing(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parents
	var expected []*v1.Parent
	for i := 0; i < 5; i++ {
		exampleParent := &v1.Parent{
			ParentName: gofakeit.Word(),
		}

		createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
			ParentName: exampleParent.ParentName,
		})
		checkValueAndError(t, createdParent, err)

		expected = append(expected, createdParent)
	}

	// Assert parent list equality
	actual, err := todoClient.GetParents(ctx, createdGrandparent.ID, nil)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual.Parents),
		"expected %d to be <= %d",
		len(expected),
		len(actual.Parents),
	)

	// Clean up
	for _, createdParent := range actual.Parents {
		err = todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID)
		assert.NoError(t, err)
	}

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("lone type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestListing(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparents
	var expected []*v1.Grandparent
	for i := 0; i < 5; i++ {
		exampleGrandparent := &v1.Grandparent{
			GrandparentName: gofakeit.Word(),
		}

		createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
			GrandparentName: exampleGrandparent.GrandparentName,
		})
		checkValueAndError(t, createdGrandparent, err)

		expected = append(expected, createdGrandparent)
	}

	// Assert grandparent list equality
	actual, err := todoClient.GetGrandparents(ctx, nil)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual.Grandparents),
		"expected %d to be <= %d",
		len(expected),
		len(actual.Grandparents),
	)

	// Clean up
	for _, createdGrandparent := range actual.Grandparents {
		err = todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID)
		assert.NoError(t, err)
	}
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Attempt to fetch nonexistent child
	_, err = todoClient.GetChild(ctx, createdGrandparent.ID, createdParent.ID, nonexistentID)
	assert.Error(t, err)

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Attempt to fetch nonexistent parent
	_, err = todoClient.GetParent(ctx, createdGrandparent.ID, nonexistentID)
	assert.Error(t, err)

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("lone type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Attempt to fetch nonexistent grandparent
	_, err = todoClient.GetGrandparent(ctx, nonexistentID)
	assert.Error(t, err)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestReadingShouldBeReadable(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldBeReadable(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create child
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	checkValueAndError(t, createdChild, err)

	// Fetch child
	actual, err := todoClient.GetChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID)
	checkValueAndError(t, actual, err)

	// Assert child equality
	checkChildEquality(t, createdChild, actual)

	// Clean up child
	assert.NoError(t, todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID))

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("single dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldBeReadable(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Fetch parent
	actual, err := todoClient.GetParent(ctx, createdGrandparent.ID, createdParent.ID)
	checkValueAndError(t, actual, err)

	// Assert parent equality
	checkParentEquality(t, createdParent, actual)

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("lone type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestReadingShouldBeReadable(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Fetch grandparent
	actual, err := todoClient.GetGrandparent(ctx, createdGrandparent.ID)
	checkValueAndError(t, actual, err)

	// Assert grandparent equality
	checkGrandparentEquality(t, createdGrandparent, actual)

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	assert.Error(t, todoClient.UpdateChild(ctx, createdGrandparent.ID, &v1.Child{ID: nonexistentID}))

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("with only one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	assert.Error(t, todoClient.UpdateParent(ctx, &v1.Parent{ID: nonexistentID}))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("for one type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	assert.Error(t, todoClient.UpdateGrandparent(ctx, &v1.Grandparent{ID: nonexistentID}))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestUpdatingShouldBeUpdateable(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldBeUpdatable(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create child
	expected := &v1.Child{
		ChildName: gofakeit.Word(),
	}
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	checkValueAndError(t, createdChild, err)

	// Change child
	createdChild.Update(expected.ToInput())
	err = todoClient.UpdateChild(ctx, createdGrandparent.ID, createdChild)
	assert.NoError(t, err)

	// Fetch child
	actual, err := todoClient.GetChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID)
	checkValueAndError(t, actual, err)

	// Assert child equality
	checkChildEquality(t, expected, actual)
	assert.NotNil(t, actual.UpdatedOn)

	// Clean up child
	assert.NoError(t, todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID))

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("with only one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldBeUpdatable(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	expected := &v1.Parent{
		ParentName: gofakeit.Word(),
	}
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Change parent
	createdParent.Update(expected.ToInput())
	err = todoClient.UpdateParent(ctx, createdParent)
	assert.NoError(t, err)

	// Fetch parent
	actual, err := todoClient.GetParent(ctx, createdGrandparent.ID, createdParent.ID)
	checkValueAndError(t, actual, err)

	// Assert parent equality
	checkParentEquality(t, expected, actual)
	assert.NotNil(t, actual.UpdatedOn)

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("for one type", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestUpdatingShouldBeUpdatable(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	expected := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Change grandparent
	createdGrandparent.Update(expected.ToInput())
	err = todoClient.UpdateGrandparent(ctx, createdGrandparent)
	assert.NoError(t, err)

	// Fetch grandparent
	actual, err := todoClient.GetGrandparent(ctx, createdGrandparent.ID)
	checkValueAndError(t, actual, err)

	// Assert grandparent equality
	checkGrandparentEquality(t, expected, actual)
	assert.NotNil(t, actual.UpdatedOn)

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDeletingShouldBeAbleToBeDeleted(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")
		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildTestDeletingShouldBeAbleToBeDeleted(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	assert "github.com/stretchr/testify/assert"
	trace "go.opencensus.io/trace"
	v1 "models/v1"
)

func doSomething() {
	tctx := context.Background()
	ctx, span := trace.StartSpan(tctx, t.Name())
	defer span.End()

	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := todoClient.CreateGrandparent(ctx, &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	checkValueAndError(t, createdGrandparent, err)

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := todoClient.CreateParent(ctx, createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	checkValueAndError(t, createdParent, err)

	// Create child
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := todoClient.CreateChild(ctx, createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	checkValueAndError(t, createdChild, err)

	// Clean up child
	assert.NoError(t, todoClient.ArchiveChild(ctx, createdGrandparent.ID, createdParent.ID, createdChild.ID))

	// Clean up parent
	assert.NoError(t, todoClient.ArchiveParent(ctx, createdGrandparent.ID, createdParent.ID))

	// Clean up grandparent
	assert.NoError(t, todoClient.ArchiveGrandparent(ctx, createdGrandparent.ID))
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
