package frontendv1

import (
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(proj))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/dev/null`

		assert.Error(t, RenderPackage(proj))
	})
}

func Test_packageDotJSON(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		actual := packageDotJSON()

		assert.Equal(t, pdjson, actual, "expected and actual do not match")
	})
}
