package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_utilsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := utilsDotGo(proj)

		expected := `
package example

import (
	"fmt"
	"regexp"
)

var (
	belongsToUserWithMandatedRestrictionRegexp    = regexp.MustCompile(` + "`" + `\+belongsToUser:\d+` + "`" + `)
	belongsToUserWithoutMandatedRestrictionRegexp = regexp.MustCompile(` + "`" + `belongsToUser:\d+` + "`" + `)
)

// ensureQueryIsRestrictedToUser takes a query and userID and ensures that query
// asks that results be restricted to a given user.
func ensureQueryIsRestrictedToUser(query string, userID uint64) string {
	switch {
	case belongsToUserWithMandatedRestrictionRegexp.MatchString(query):
		return query
	case belongsToUserWithoutMandatedRestrictionRegexp.MatchString(query):
		query = belongsToUserWithoutMandatedRestrictionRegexp.ReplaceAllString(query, fmt.Sprintf("+belongsToUser:%d", userID))
	case !belongsToUserWithMandatedRestrictionRegexp.MatchString(query):
		query = fmt.Sprintf("%s +belongsToUser:%d", query, userID)
	}

	return query
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUtilsVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildUtilsVarDeclarations()

		expected := `
package example

import (
	"fmt"
	"regexp"
)

var (
	belongsToUserWithMandatedRestrictionRegexp    = regexp.MustCompile(` + "`" + `\+belongsToUser:\d+` + "`" + `)
	belongsToUserWithoutMandatedRestrictionRegexp = regexp.MustCompile(` + "`" + `belongsToUser:\d+` + "`" + `)
)

// ensureQueryIsRestrictedToUser takes a query and userID and ensures that query
// asks that results be restricted to a given user.
func ensureQueryIsRestrictedToUser(query string, userID uint64) string {
	switch {
	case belongsToUserWithMandatedRestrictionRegexp.MatchString(query):
		return query
	case belongsToUserWithoutMandatedRestrictionRegexp.MatchString(query):
		query = belongsToUserWithoutMandatedRestrictionRegexp.ReplaceAllString(query, fmt.Sprintf("+belongsToUser:%d", userID))
	case !belongsToUserWithMandatedRestrictionRegexp.MatchString(query):
		query = fmt.Sprintf("%s +belongsToUser:%d", query, userID)
	}

	return query
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
