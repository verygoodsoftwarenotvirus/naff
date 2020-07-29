package client

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := docDotGo()

		expected := `
/*
Package client provides an HTTP client that can communicate with and interpret the responses
of an instance of the todo service.
*/
package client

import ()
`
		b := bytes.NewBufferString("\n")
		require.NoError(t, x.Render(b))

		actual := b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
