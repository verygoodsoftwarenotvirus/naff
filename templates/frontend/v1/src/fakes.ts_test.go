package frontendsrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fakesDotTS(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := fakesDotTSFile
		actual := fakesDotTS()

		assert.Equal(t, expected, actual, "expected and actual do not match")
	})
}
