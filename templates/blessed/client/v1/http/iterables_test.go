package client

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func Test_buildParamsForMethodThatHandlesAnInstanceOfADataType(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, cherry)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID, bananaID, cherryID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation with fewer dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, banana)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID, bananaID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation without dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
			BelongsToNobody: true,
		}
		proj := &models.Project{
			DataTypes: []models.DataType{apple},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, apple)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildGetSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildGetSomethingRequestFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildGetSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildGetSomethingFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildGetSomethingsRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildGetSomethingsRequestFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildGetListOfSomethingsFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildGetListOfSomethingsFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildCreateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildCreateSomethingRequestFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildCreateSomethingFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildUpdateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildUpdateSomethingRequestFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildUpdateSomethingFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildArchiveSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildArchiveSomethingRequestFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildArchiveSomethingFuncDecl(proj, cherry)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"net/http"
	"strconv"
)

// BuildGetCherryRequest builds an HTTP request for fetching a cherry
func (c *V1Client) BuildGetCherryRequest(ctx context.Context, appleID, bananaID, cherryID uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, applesBasePath, strconv.FormatUint(appleID, 10), bananasBasePath, strconv.FormatUint(bananaID, 10), cherriesBasePath, strconv.FormatUint(cherryID, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
