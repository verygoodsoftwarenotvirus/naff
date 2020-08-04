package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_actionsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := actionsDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	http1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"math/rand"
	"net/http"
	"time"
)

var (
	// ErrUnavailableYet is a sentinel error value.
	ErrUnavailableYet = errors.New("can't do this yet")
)

type (
	// actionFunc represents a thing you can do.
	actionFunc func() (*http.Request, error)

	// Action is a wrapper struct around some important values.
	Action struct {
		Action actionFunc
		Weight int
		Name   string
	}
)

// RandomAction takes a client and returns a closure which is an action.
func RandomAction(c *http1.V1Client) *Action {
	allActions := map[string]*Action{
		"GetHealthCheck": {
			Name: "GetHealthCheck",
			Action: func() (*http.Request, error) {
				ctx := context.Background()
				return c.BuildHealthCheckRequest(ctx)
			},
			Weight: 100,
		},
		"CreateUser": {
			Name: "CreateUser",
			Action: func() (*http.Request, error) {
				ctx := context.Background()
				ui := fake.BuildFakeUserCreationInput()
				return c.BuildCreateUserRequest(ctx, ui)
			},
			Weight: 100,
		},
	}

	for k, v := range buildItemActions(c) {
		allActions[k] = v
	}

	for k, v := range buildWebhookActions(c) {
		allActions[k] = v
	}

	for k, v := range buildOAuth2ClientActions(c) {
		allActions[k] = v
	}

	totalWeight := 0
	for _, rb := range allActions {
		totalWeight += rb.Weight
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)

	for _, rb := range allActions {
		r -= rb.Weight
		if r <= 0 {
			return rb
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildVarDefs()

		expected := `
package example

import (
	"errors"
)

var (
	// ErrUnavailableYet is a sentinel error value.
	ErrUnavailableYet = errors.New("can't do this yet")
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTypeDefs()

		expected := `
package example

import (
	"net/http"
)

type (
	// actionFunc represents a thing you can do.
	actionFunc func() (*http.Request, error)

	// Action is a wrapper struct around some important values.
	Action struct {
		Action actionFunc
		Weight int
		Name   string
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandomAction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildRandomAction(proj)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"math/rand"
	http1 "net/http"
	"time"
)

// RandomAction takes a client and returns a closure which is an action.
func RandomAction(c *http.V1Client) *Action {
	allActions := map[string]*Action{
		"GetHealthCheck": {
			Name: "GetHealthCheck",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				return c.BuildHealthCheckRequest(ctx)
			},
			Weight: 100,
		},
		"CreateUser": {
			Name: "CreateUser",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				ui := fake.BuildFakeUserCreationInput()
				return c.BuildCreateUserRequest(ctx, ui)
			},
			Weight: 100,
		},
	}

	for k, v := range buildItemActions(c) {
		allActions[k] = v
	}

	for k, v := range buildWebhookActions(c) {
		allActions[k] = v
	}

	for k, v := range buildOAuth2ClientActions(c) {
		allActions[k] = v
	}

	totalWeight := 0
	for _, rb := range allActions {
		totalWeight += rb.Weight
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)

	for _, rb := range allActions {
		r -= rb.Weight
		if r <= 0 {
			return rb
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
