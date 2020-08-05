package project

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderProject(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = filepath.Join(os.TempDir(), "things", "stuff")

		assert.NoError(t, RenderProject(proj))
	})

	//T.Run("with invalid output directory", func(t *testing.T) {
	//	proj := testprojects.BuildTodoApp()
	//	proj.OutputPath = `/dev/null`
	//
	//	assert.Panics(t, func() { RenderProject(proj) })
	//})
}
