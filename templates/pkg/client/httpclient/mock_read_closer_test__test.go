package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockReadCloserTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mockReadCloserTestDotGo(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	"io"
)

var _ io.ReadCloser = (*ReadCloser)(nil)

// ReadCloser is a mock io.ReadCloser for testing purposes.
type ReadCloser struct {
	mock.Mock
}

// newMockReadCloser returns a new mock io.ReadCloser.
func newMockReadCloser() *ReadCloser {
	return &ReadCloser{}
}

// ReadHandler implements the ReadHandler part of our ReadCloser.
func (m *ReadCloser) Read(b []byte) (i int, err error) {
	retVals := m.Called(b)
	return retVals.Int(0), retVals.Error(1)
}

// Close implements the Closer part of our ReadCloser.
func (m *ReadCloser) Close() (err error) {
	return m.Called().Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserInterfaceAssurance(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockReadCloserInterfaceAssurance()

		expected := `
package example

import (
	"io"
)

var _ io.ReadCloser = (*ReadCloser)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockReadCloserDecl()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// ReadCloser is a mock io.ReadCloser for testing purposes.
type ReadCloser struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewMockReadCloser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildNewMockReadCloser()

		expected := `
package example

import ()

// newMockReadCloser returns a new mock io.ReadCloser.
func newMockReadCloser() *ReadCloser {
	return &ReadCloser{}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockReadCloserReadHandler()

		expected := `
package example

import ()

// ReadHandler implements the ReadHandler part of our ReadCloser.
func (m *ReadCloser) Read(b []byte) (i int, err error) {
	retVals := m.Called(b)
	return retVals.Int(0), retVals.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserClose(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockReadCloserClose()

		expected := `
package example

import ()

// Close implements the Closer part of our ReadCloser.
func (m *ReadCloser) Close() (err error) {
	return m.Called().Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
