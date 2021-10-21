package frontendsrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_queryFilterDotTS(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := queryFilterDotTSFile
		actual := queryFilterDotTS()

		assert.Equal(t, expected, actual, "expected and actual do not match")
	})
}
