package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_coverageTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := coverageTestDotGo(proj)

		expected := `
package example

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestRunMain(_ *testing.T) {
	// This test is built specifically to capture the coverage that the integration
	// tests exhibit. We run the main function (i.e. a production server)
	// on an independent goroutine and sleep for long enough that the integration
	// tests can run, then we quit.
	d, err := time.ParseDuration(os.Getenv("RUNTIME_DURATION"))
	if err != nil {
		log.Fatal(err)
	}

	go main()

	time.Sleep(d)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestMain()

		expected := `
package example

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestRunMain(_ *testing.T) {
	// This test is built specifically to capture the coverage that the integration
	// tests exhibit. We run the main function (i.e. a production server)
	// on an independent goroutine and sleep for long enough that the integration
	// tests can run, then we quit.
	d, err := time.ParseDuration(os.Getenv("RUNTIME_DURATION"))
	if err != nil {
		log.Fatal(err)
	}

	go main()

	time.Sleep(d)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
