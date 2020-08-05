package frontendsrc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stateDotTS(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := stateDotTSFile
		actual := stateDotTS()

		assert.Equal(t, expected, actual, "expected and actual do not match")
	})
}
