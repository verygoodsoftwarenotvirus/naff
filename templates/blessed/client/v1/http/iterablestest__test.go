package client

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"
)

func Test_buildTestV1Client_BuildGetSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildTestV1Client_BuildGetSomethingRequest(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildGetChildRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		expectedID := uint64(1)
		actual, err := c.BuildGetChildRequest(ctx, grandparentID, parentID, childID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", expectedID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestV1Client_GetSomething(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildTestV1Client_GetSomething(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "models/v1"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_GetChild(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		grandparent := &v1.Grandparent{ID: 1}
		parent := &v1.Parent{
			ID:                   1,
			BelongsToGrandparent: grandparent.ID,
		}
		child := &v1.Child{
			ID:              1,
			BelongsToParent: parent.ID,
		}

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(child.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/grandparents/%d/parents/%d/children/%d", grandparent.ID, parent.ID, child.ID), "expected and actual path don't match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(child))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetChild(ctx, grandparent.ID, parent.ID, child.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, child, actual)
	})
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestV1Client_BuildGetListOfSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildTestV1Client_BuildGetListOfSomethingRequest(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "models/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildGetChildrenRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		filter := (*v1.QueryFilter)(nil)
		grandparentID := uint64(1)
		parentID := uint64(1)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetChildrenRequest(ctx, grandparentID, parentID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
