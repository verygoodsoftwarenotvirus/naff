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

var (
	a = models.DataType{
		Name: wordsmith.FromSingularPascalCase("Grandparent"),
		Fields: []models.DataField{
			{
				Name: wordsmith.FromSingularPascalCase("GrandparentName"),
				Type: "string",
			},
		},
	}
	b = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Parent"),
		BelongsToStruct: a.Name,
		Fields: []models.DataField{
			{
				Name: wordsmith.FromSingularPascalCase("ParentName"),
				Type: "string",
			},
		},
	}
	c = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Child"),
		BelongsToStruct: b.Name,
		Fields: []models.DataField{
			{
				Name: wordsmith.FromSingularPascalCase("ChildName"),
				Type: "string",
			},
		},
	}
)

func Test_buildParamsForMethodThatHandlesAnInstanceOfADataType(T *testing.T) {
	T.Parallel()

	proj := &models.Project{
		DataTypes: []models.DataType{a, b, c},
	}

	T.Run("normal operation with dependencies", func(t *testing.T) {
		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, c, false)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, grandparentID, parentID, childID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation with fewer dependencies", func(t *testing.T) {
		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, b, false)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, grandparentID, parentID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation without dependencies", func(t *testing.T) {
		gp := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Grandparent"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("Name"),
					Type: "string",
				},
			},
			BelongsToNobody: true,
		}
		proj := &models.Project{
			DataTypes: []models.DataType{gp},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, a, false)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, grandparentID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildParamsForMethodThatIncludesItsOwnTypeInItsParams(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(
				buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, c, false)...,
			).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
)

func doSomething(ctx context.Context, grandparentID uint64, child *v1.Child) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation with fewer dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, b, false)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
)

func doSomething(ctx context.Context, parent *v1.Parent) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation without dependencies", func(t *testing.T) {
		gp := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Grandparent"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("Name"),
					Type: "string",
				},
			},
			BelongsToNobody: true,
		}
		proj := &models.Project{
			DataTypes: []models.DataType{gp},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj, a, false)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
)

func doSomething(ctx context.Context, grandparent *v1.Grandparent) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildGetSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildGetSomethingRequestFuncDecl(proj, c)...,
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

// BuildGetChildRequest builds an HTTP request for fetching a child
func (c *V1Client) BuildGetChildRequest(ctx context.Context, grandparentID, parentID, childID uint64) (*http.Request, error) {
	uri := c.BuildURL(
		nil,
		grandparentsBasePath,
		strconv.FormatUint(grandparentID, 10),
		parentsBasePath,
		strconv.FormatUint(parentID, 10),
		childrenBasePath,
		strconv.FormatUint(childID, 10),
	)

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
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildGetSomethingFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
	v1 "models/v1"
)

// GetChild retrieves a child
func (c *V1Client) GetChild(ctx context.Context, grandparentID, parentID, childID uint64) (child *v1.Child, err error) {
	req, err := c.BuildGetChildRequest(ctx, grandparentID, parentID, childID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &child); retrieveErr != nil {
		return nil, retrieveErr
	}

	return child, nil
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildGetSomethingsRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildGetSomethingsRequestFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
	"net/http"
	"strconv"
)

// BuildGetChildrenRequest builds an HTTP request for fetching children
func (c *V1Client) BuildGetChildrenRequest(ctx context.Context, grandparentID, parentID uint64, filter *v1.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(
		filter.ToValues(),
		grandparentsBasePath,
		strconv.FormatUint(grandparentID, 10),
		parentsBasePath,
		strconv.FormatUint(parentID, 10),
		childrenBasePath,
	)

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
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildGetListOfSomethingsFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
	v1 "models/v1"
)

// GetChildren retrieves a list of children
func (c *V1Client) GetChildren(ctx context.Context, grandparentID, parentID uint64, filter *v1.QueryFilter) (children *v1.ChildList, err error) {
	req, err := c.BuildGetChildrenRequest(ctx, grandparentID, parentID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &children); retrieveErr != nil {
		return nil, retrieveErr
	}

	return children, nil
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildCreateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildCreateSomethingRequestFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
	"net/http"
	"strconv"
)

// BuildCreateChildRequest builds an HTTP request for creating a child
func (c *V1Client) BuildCreateChildRequest(ctx context.Context, grandparentID, parentID uint64, input *v1.ChildCreationInput) (*http.Request, error) {
	uri := c.BuildURL(
		nil,
		grandparentsBasePath,
		strconv.FormatUint(grandparentID, 10),
		parentsBasePath,
		strconv.FormatUint(parentID, 10),
		childrenBasePath,
	)

	return c.buildDataRequest(http.MethodPost, uri, input)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildCreateSomethingFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
	v1 "models/v1"
)

// CreateChild creates a child
func (c *V1Client) CreateChild(ctx context.Context, grandparentID, parentID uint64, input *v1.ChildCreationInput) (child *v1.Child, err error) {
	req, err := c.BuildCreateChildRequest(ctx, grandparentID, parentID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &child)
	return child, err
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildUpdateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildUpdateSomethingRequestFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	v1 "models/v1"
	"net/http"
	"strconv"
)

// BuildUpdateChildRequest builds an HTTP request for updating a child
func (c *V1Client) BuildUpdateChildRequest(ctx context.Context, grandparentID uint64, child *v1.Child) (*http.Request, error) {
	uri := c.BuildURL(
		nil,
		grandparentsBasePath,
		strconv.FormatUint(grandparentID, 10),
		parentsBasePath,
		strconv.FormatUint(child.BelongsToParent, 10),
		childrenBasePath,
		strconv.FormatUint(child.ID, 10),
	)

	return c.buildDataRequest(http.MethodPut, uri, child)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildUpdateSomethingFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
	v1 "models/v1"
)

// UpdateChild updates a child
func (c *V1Client) UpdateChild(ctx context.Context, grandparentID uint64, child *v1.Child) error {
	req, err := c.BuildUpdateChildRequest(ctx, grandparentID, child)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &child)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildBuildArchiveSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildBuildArchiveSomethingRequestFuncDecl(proj, c)...,
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

// BuildArchiveChildRequest builds an HTTP request for updating a child
func (c *V1Client) BuildArchiveChildRequest(ctx context.Context, grandparentID, parentID, childID uint64) (*http.Request, error) {
	uri := c.BuildURL(
		nil,
		grandparentsBasePath,
		strconv.FormatUint(grandparentID, 10),
		parentsBasePath,
		strconv.FormatUint(parentID, 10),
		childrenBasePath,
		strconv.FormatUint(childID, 10),
	)

	return http.NewRequest(http.MethodDelete, uri, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func TestBuildArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildArchiveSomethingFuncDecl(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
)

// ArchiveChild archives a child
func (c *V1Client) ArchiveChild(ctx context.Context, grandparentID, parentID, childID uint64) error {
	req, err := c.BuildArchiveChildRequest(ctx, grandparentID, parentID, childID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
