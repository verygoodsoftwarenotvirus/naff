package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientTestDotGo(proj)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_GetID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			ClientID: expected,
		}
		assert.Equal(t, expected, oac.GetID())
	})
}

func TestOAuth2Client_GetSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			ClientSecret: expected,
		}
		assert.Equal(t, expected, oac.GetSecret())
	})
}

func TestOAuth2Client_GetDomain(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			RedirectURI: expected,
		}
		assert.Equal(t, expected, oac.GetDomain())
	})
}

func TestOAuth2Client_GetUserID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := uint64(123)
		expected := fmt.Sprintf("%d", expectation)
		oac := &OAuth2Client{
			BelongsToUser: expectation,
		}
		assert.Equal(t, expected, oac.GetUserID())
	})
}

func TestOAuth2Client_HasScope(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		oac := &OAuth2Client{
			Scopes: []string{"things", "and", "stuff"},
		}

		assert.True(t, oac.HasScope(oac.Scopes[0]))
		assert.False(t, oac.HasScope("blah"))
		assert.False(t, oac.HasScope(""))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestOAuth2Client_GetID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestOAuth2Client_GetID()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_GetID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			ClientID: expected,
		}
		assert.Equal(t, expected, oac.GetID())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestOAuth2Client_GetSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestOAuth2Client_GetSecret()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_GetSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			ClientSecret: expected,
		}
		assert.Equal(t, expected, oac.GetSecret())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestOAuth2Client_GetDomain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestOAuth2Client_GetDomain()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_GetDomain(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expected := "123"
		oac := &OAuth2Client{
			RedirectURI: expected,
		}
		assert.Equal(t, expected, oac.GetDomain())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestOAuth2Client_GetUserID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestOAuth2Client_GetUserID()

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_GetUserID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := uint64(123)
		expected := fmt.Sprintf("%d", expectation)
		oac := &OAuth2Client{
			BelongsToUser: expectation,
		}
		assert.Equal(t, expected, oac.GetUserID())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestOAuth2Client_HasScope(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestOAuth2Client_HasScope()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOAuth2Client_HasScope(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		oac := &OAuth2Client{
			Scopes: []string{"things", "and", "stuff"},
		}

		assert.True(t, oac.HasScope(oac.Scopes[0]))
		assert.False(t, oac.HasScope("blah"))
		assert.False(t, oac.HasScope(""))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
