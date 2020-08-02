package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_initDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := initDotGo(proj)

		expected := `
package example

import (
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"strings"
	"time"
)

var urlToUse string

const (
	seleniumHubAddr = "http://selenium-hub:4444/wd/hub"
)

func init() {
	urlToUse = testutil.DetermineServiceURL()

	logger := zerolog.NewZeroLogger()
	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	// NOTE: this is sad, but also the only thing that consistently works
	// see above for my vain attempts at a real solution to this problem.
	time.Sleep(10 * time.Second)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInitVarDefs()

		expected := `
package example

import ()

var urlToUse string
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInitConstDefs()

		expected := `
package example

import ()

const (
	seleniumHubAddr = "http://selenium-hub:4444/wd/hub"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInitInit(proj)

		expected := `
package example

import (
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"strings"
	"time"
)

func init() {
	urlToUse = testutil.DetermineServiceURL()

	logger := zerolog.NewZeroLogger()
	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	// NOTE: this is sad, but also the only thing that consistently works
	// see above for my vain attempts at a real solution to this problem.
	time.Sleep(10 * time.Second)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
