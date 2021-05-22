package authentication

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireTestDotGo(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	"testing"
)

func TestProvideWebsocketAuthFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideWebsocketAuthFunc(buildTestService(t)))
	})
}

func TestProvideOAuth2ClientValidator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideOAuth2ClientValidator(&oauth2clients.Service{}))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideWebsocketAuthFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideWebsocketAuthFunc()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestProvideWebsocketAuthFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideWebsocketAuthFunc(buildTestService(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideOAuth2ClientValidator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestProvideOAuth2ClientValidator(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	"testing"
)

func TestProvideOAuth2ClientValidator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideOAuth2ClientValidator(&oauth2clients.Service{}))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
