package composefiles

import (
	"os"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()
		project.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(project))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/dev/null`

		assert.Error(t, RenderPackage(proj))
	})
}
