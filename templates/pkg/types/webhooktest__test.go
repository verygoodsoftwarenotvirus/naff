package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhookTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhookTestDotGo(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"net/http"
	"testing"
)

func TestWebhook_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleInput := &WebhookUpdateInput{
			Name:        "whatever",
			ContentType: "application/xml",
			URL:         "https://blah.verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPatch,
			Events:      []string{"more_things"},
			DataTypes:   []string{"new_stuff"},
			Topics:      []string{"blah-blah"},
		}

		actual := &Webhook{
			Name:        "something_else",
			ContentType: "application/json",
			URL:         "https://verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPost,
			Events:      []string{"things"},
			DataTypes:   []string{"stuff"},
			Topics:      []string{"blah"},
		}
		expected := &Webhook{
			Name:        exampleInput.Name,
			ContentType: "application/xml",
			URL:         "https://blah.verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPatch,
			Events:      []string{"more_things"},
			DataTypes:   []string{"new_stuff"},
			Topics:      []string{"blah-blah"},
		}

		actual.Update(exampleInput)
		assert.Equal(t, expected, actual)
	})
}

func TestWebhook_ToListener(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		w := &Webhook{}
		w.ToListener(noop.ProvideNoopLogger())
	})
}

func Test_buildErrorLogFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		w := &Webhook{}
		actual := buildErrorLogFunc(w, noop.ProvideNoopLogger())
		actual(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhook_Update(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestWebhook_Update()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWebhook_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleInput := &WebhookUpdateInput{
			Name:        "whatever",
			ContentType: "application/xml",
			URL:         "https://blah.verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPatch,
			Events:      []string{"more_things"},
			DataTypes:   []string{"new_stuff"},
			Topics:      []string{"blah-blah"},
		}

		actual := &Webhook{
			Name:        "something_else",
			ContentType: "application/json",
			URL:         "https://verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPost,
			Events:      []string{"things"},
			DataTypes:   []string{"stuff"},
			Topics:      []string{"blah"},
		}
		expected := &Webhook{
			Name:        exampleInput.Name,
			ContentType: "application/xml",
			URL:         "https://blah.verygoodsoftwarenotvirus.ru",
			Method:      http.MethodPatch,
			Events:      []string{"more_things"},
			DataTypes:   []string{"new_stuff"},
			Topics:      []string{"blah-blah"},
		}

		actual.Update(exampleInput)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhook_ToListener(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestWebhook_ToListener()

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestWebhook_ToListener(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		w := &Webhook{}
		w.ToListener(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_buildErrorLogFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_buildErrorLogFunc()

		expected := `
package example

import (
	"errors"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func Test_buildErrorLogFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		w := &Webhook{}
		actual := buildErrorLogFunc(w, noop.ProvideNoopLogger())
		actual(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
