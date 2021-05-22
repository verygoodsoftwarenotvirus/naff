package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_serviceDataEventsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := serviceDataEventsDotGo(proj)

		expected := `
package example

import ()

// ServiceDataEvent is a simple string alias.
type ServiceDataEvent string

const (
	// Create represents a create event.
	Create ServiceDataEvent = "create"
	// Update represents an update event.
	Update ServiceDataEvent = "update"
	// Archive represents an archive event.
	Archive ServiceDataEvent = "archive"
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceDataEvent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceDataEvent()

		expected := `
package example

import ()

// ServiceDataEvent is a simple string alias.
type ServiceDataEvent string
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceDataEventsConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceDataEventsConstantDefs()

		expected := `
package example

import ()

const (
	// Create represents a create event.
	Create ServiceDataEvent = "create"
	// Update represents an update event.
	Update ServiceDataEvent = "update"
	// Archive represents an archive event.
	Archive ServiceDataEvent = "archive"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
