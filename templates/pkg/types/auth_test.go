package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := authDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"encoding/gob"
)

const (
	// SessionInfoKey is the non-string type we use for referencing SessionInfo structs
	SessionInfoKey ContextKey = "session_info"
)

func init() {
	gob.Register(&SessionInfo{})
}

type (
	// SessionInfo represents what we encode in our authentication cookies.
	SessionInfo struct {
		UserID      uint64 ` + "`" + `json:"-"` + "`" + `
		UserIsAdmin bool   ` + "`" + `json:"-"` + "`" + `
	}

	// StatusResponse is what we encode when the frontend wants to check auth status
	StatusResponse struct {
		Authenticated bool ` + "`" + `json:"isAuthenticated"` + "`" + `
		IsAdmin       bool ` + "`" + `json:"isAdmin"` + "`" + `
	}
)

// ToBytes returns the gob encoded session info
func (i *SessionInfo) ToBytes() []byte {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(i); err != nil {
		panic(err)
	}

	return b.Bytes()
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthConstantDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAuthConstantDefinitions()

		expected := `
package example

import ()

const (
	// SessionInfoKey is the non-string type we use for referencing SessionInfo structs
	SessionInfoKey ContextKey = "session_info"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAuthInit()

		expected := `
package example

import (
	"encoding/gob"
)

func init() {
	gob.Register(&SessionInfo{})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAuthTypeDefinitions()

		expected := `
package example

import ()

type (
	// SessionInfo represents what we encode in our authentication cookies.
	SessionInfo struct {
		UserID      uint64 ` + "`" + `json:"-"` + "`" + `
		UserIsAdmin bool   ` + "`" + `json:"-"` + "`" + `
	}

	// StatusResponse is what we encode when the frontend wants to check auth status
	StatusResponse struct {
		Authenticated bool ` + "`" + `json:"isAuthenticated"` + "`" + `
		IsAdmin       bool ` + "`" + `json:"isAdmin"` + "`" + `
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthSessionInfoToBytes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAuthSessionInfoToBytes()

		expected := `
package example

import (
	"bytes"
	"encoding/gob"
)

// ToBytes returns the gob encoded session info
func (i *SessionInfo) ToBytes() []byte {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(i); err != nil {
		panic(err)
	}

	return b.Bytes()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
