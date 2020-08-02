package v1

import (
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func Test_jsonTag(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		expected := map[string]string{
			"json": "fart",
		}
		actual := jsonTag("fart")

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(proj))
	})
}
