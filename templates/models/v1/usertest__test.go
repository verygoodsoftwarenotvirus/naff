package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_userTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := userTestDotGo(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		actual := User{
			Username:        "old_username",
			HashedPassword:  "hashed_pass",
			TwoFactorSecret: "two factor secret",
		}
		exampleInput := User{
			Username:        "new_username",
			HashedPassword:  "updated_hashed_pass",
			TwoFactorSecret: "new fancy secret",
		}

		actual.Update(&exampleInput)
		assert.Equal(t, exampleInput, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestUser_Update(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestUser_Update()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		actual := User{
			Username:        "old_username",
			HashedPassword:  "hashed_pass",
			TwoFactorSecret: "two factor secret",
		}
		exampleInput := User{
			Username:        "new_username",
			HashedPassword:  "updated_hashed_pass",
			TwoFactorSecret: "new fancy secret",
		}

		actual.Update(&exampleInput)
		assert.Equal(t, exampleInput, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
