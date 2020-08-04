package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(proj))
	})

	//	T.Run("with invalid output directory", func(t *testing.T) {
	//		t.Parallel()
	//
	//		proj := testprojects.BuildTodoApp()
	//		proj.OutputPath = `/\0/\0/\0`
	//
	//		assert.Error(t, RenderPackage(proj))
	//	})
}
