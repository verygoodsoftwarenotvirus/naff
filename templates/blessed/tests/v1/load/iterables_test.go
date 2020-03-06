package load

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

func Test_buildRequisiteCreationCode(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildRequisiteCreationCode(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	v1 "models/v1"
)

func doSomething() {
	// Create grandparent
	exampleGrandparent := &v1.Grandparent{
		GrandparentName: gofakeit.Word(),
	}

	createdGrandparent, err := c.CreateGrandparent(context.Background(), &v1.GrandparentCreationInput{
		GrandparentName: exampleGrandparent.GrandparentName,
	})
	if err != nil {
		return nil, err
	}

	// Create parent
	exampleParent := &v1.Parent{
		ParentName: gofakeit.Word(),
	}

	createdParent, err := c.CreateParent(context.Background(), createdGrandparent.ID, &v1.ParentCreationInput{
		ParentName: exampleParent.ParentName,
	})
	if err != nil {
		return nil, err
	}

	// Create child
	exampleChild := &v1.Child{
		ChildName: gofakeit.Word(),
	}

	createdChild, err := c.CreateChild(context.Background(), createdGrandparent.ID, createdParent.ID, &v1.ChildCreationInput{
		ChildName: exampleChild.ChildName,
	})
	if err != nil {
		return nil, err
	}

}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildRandomActionMap(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildRandomActionMap(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	v1 "models/v1"
	http1 "net/http"
	model "tests/v1/testutil/rand/model"
)

func buildChildActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateChild": {
			Name: "CreateChild",
			Action: func() (*http1.Request, error) {
				// Create grandparent
				exampleGrandparent := &v1.Grandparent{
					GrandparentName: gofakeit.Word(),
				}

				createdGrandparent, err := c.CreateGrandparent(context.Background(), &v1.GrandparentCreationInput{
					GrandparentName: exampleGrandparent.GrandparentName,
				})
				if err != nil {
					return nil, err
				}

				// Create parent
				exampleParent := &v1.Parent{
					ParentName: gofakeit.Word(),
				}

				createdParent, err := c.CreateParent(context.Background(), createdGrandparent.ID, &v1.ParentCreationInput{
					ParentName: exampleParent.ParentName,
				})
				if err != nil {
					return nil, err
				}

				req, err := c.BuildCreateChildRequest(context.Background(), createdGrandparent.ID, createdParent.ID, model.RandomChildCreationInput())

				return req, err
			},
			Weight: 100,
		},
		"GetChild": {
			Name: "GetChild",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)
				parent := fetchRandomParent(c, grandparent.ID)

				if randomChild := fetchRandomChild(c, grandparent.ID, parent.ID); randomChild != nil {
					return c.BuildGetChildRequest(context.Background(), grandparent.ID, parent.ID, randomChild.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetChildren": {
			Name: "GetChildren",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)
				parent := fetchRandomParent(c, grandparent.ID)

				return c.BuildGetChildrenRequest(context.Background(), grandparent.ID, parent.ID, nil)
			},
			Weight: 100,
		},
		"UpdateChild": {
			Name: "UpdateChild",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)
				parent := fetchRandomParent(c, grandparent.ID)

				if randomChild := fetchRandomChild(c, grandparent.ID, parent.ID); randomChild != nil {
					randomChild.ChildName = model.RandomChildCreationInput().ChildName

					return c.BuildUpdateChildRequest(context.Background(), grandparent.ID, randomChild)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveChild": {
			Name: "ArchiveChild",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)
				parent := fetchRandomParent(c, grandparent.ID)

				if randomChild := fetchRandomChild(c, grandparent.ID, parent.ID); randomChild != nil {
					return c.BuildArchiveChildRequest(context.Background(), grandparent.ID, parent.ID, randomChild.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("with one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildRandomActionMap(proj, b)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	gofakeit "github.com/brianvoe/gofakeit"
	v1 "models/v1"
	http1 "net/http"
	model "tests/v1/testutil/rand/model"
)

func buildParentActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateParent": {
			Name: "CreateParent",
			Action: func() (*http1.Request, error) {
				// Create grandparent
				exampleGrandparent := &v1.Grandparent{
					GrandparentName: gofakeit.Word(),
				}

				createdGrandparent, err := c.CreateGrandparent(context.Background(), &v1.GrandparentCreationInput{
					GrandparentName: exampleGrandparent.GrandparentName,
				})
				if err != nil {
					return nil, err
				}

				req, err := c.BuildCreateParentRequest(context.Background(), createdGrandparent.ID, model.RandomParentCreationInput())

				return req, err
			},
			Weight: 100,
		},
		"GetParent": {
			Name: "GetParent",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)

				if randomParent := fetchRandomParent(c, grandparent.ID); randomParent != nil {
					return c.BuildGetParentRequest(context.Background(), grandparent.ID, randomParent.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetParents": {
			Name: "GetParents",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)

				return c.BuildGetParentsRequest(context.Background(), grandparent.ID, nil)
			},
			Weight: 100,
		},
		"UpdateParent": {
			Name: "UpdateParent",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)

				if randomParent := fetchRandomParent(c, grandparent.ID); randomParent != nil {
					randomParent.ParentName = model.RandomParentCreationInput().ParentName

					return c.BuildUpdateParentRequest(context.Background(), randomParent)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveParent": {
			Name: "ArchiveParent",
			Action: func() (*http1.Request, error) {
				grandparent := fetchRandomGrandparent(c)

				if randomParent := fetchRandomParent(c, grandparent.ID); randomParent != nil {
					return c.BuildArchiveParentRequest(context.Background(), grandparent.ID, randomParent.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
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
			buildRandomActionMap(proj, a)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	http1 "net/http"
	model "tests/v1/testutil/rand/model"
)

func buildGrandparentActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateGrandparent": {
			Name: "CreateGrandparent",
			Action: func() (*http1.Request, error) {
				req, err := c.BuildCreateGrandparentRequest(context.Background(), model.RandomGrandparentCreationInput())

				return req, err
			},
			Weight: 100,
		},
		"GetGrandparent": {
			Name: "GetGrandparent",
			Action: func() (*http1.Request, error) {
				if randomGrandparent := fetchRandomGrandparent(c); randomGrandparent != nil {
					return c.BuildGetGrandparentRequest(context.Background(), randomGrandparent.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetGrandparents": {
			Name: "GetGrandparents",
			Action: func() (*http1.Request, error) {
				return c.BuildGetGrandparentsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateGrandparent": {
			Name: "UpdateGrandparent",
			Action: func() (*http1.Request, error) {
				if randomGrandparent := fetchRandomGrandparent(c); randomGrandparent != nil {
					randomGrandparent.GrandparentName = model.RandomGrandparentCreationInput().GrandparentName

					return c.BuildUpdateGrandparentRequest(context.Background(), randomGrandparent)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveGrandparent": {
			Name: "ArchiveGrandparent",
			Action: func() (*http1.Request, error) {
				if randomGrandparent := fetchRandomGrandparent(c); randomGrandparent != nil {
					return c.BuildArchiveGrandparentRequest(context.Background(), randomGrandparent.ID)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildRandomDependentIDFetchers(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildRandomDependentIDFetchers(proj, c)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import ()

func doSomething() {
	grandparent := fetchRandomGrandparent(c)
	parent := fetchRandomParent(c, grandparent.ID)
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("just one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildRandomDependentIDFetchers(proj, b)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import ()

func doSomething() {
	grandparent := fetchRandomGrandparent(c)
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
				buildRandomDependentIDFetchers(proj, a)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import ()

func doSomething() {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_buildFetchRandomSomething(T *testing.T) {
	//T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			buildFetchRandomSomething(proj, c)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	"math/rand"
	v1 "models/v1"
)

// fetchRandomChild retrieves a random child from the list of available children
func fetchRandomChild(c *http.V1Client, grandparentID, parentID uint64) *v1.Child {
	childrenRes, err := c.GetChildren(context.Background(), grandparentID, parentID, nil)
	if err != nil || childrenRes == nil || len(childrenRes.Children) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(childrenRes.Children))
	return &childrenRes.Children[randIndex]
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
			buildFetchRandomSomething(proj, b)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	"math/rand"
	v1 "models/v1"
)

// fetchRandomParent retrieves a random parent from the list of available parents
func fetchRandomParent(c *http.V1Client, grandparentID uint64) *v1.Parent {
	parentsRes, err := c.GetParents(context.Background(), grandparentID, nil)
	if err != nil || parentsRes == nil || len(parentsRes.Parents) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(parentsRes.Parents))
	return &parentsRes.Parents[randIndex]
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
			buildFetchRandomSomething(proj, a)...,
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	http "client/v1/http"
	"context"
	"math/rand"
	v1 "models/v1"
)

// fetchRandomGrandparent retrieves a random grandparent from the list of available grandparents
func fetchRandomGrandparent(c *http.V1Client) *v1.Grandparent {
	grandparentsRes, err := c.GetGrandparents(context.Background(), nil)
	if err != nil || grandparentsRes == nil || len(grandparentsRes.Grandparents) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(grandparentsRes.Grandparents))
	return &grandparentsRes.Grandparents[randIndex]
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
