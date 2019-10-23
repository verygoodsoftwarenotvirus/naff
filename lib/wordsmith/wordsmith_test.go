package wordsmith

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr(T *testing.T) {

	type test struct {
		expected, actual string
	}

	T.Run("Item", func(t *testing.T) {
		s := FromSingularPascalCase("Item")

		expectationsMap := map[string]test{
			"Singular":                {expected: "Item", actual: s.Singular()},
			"Plural":                  {expected: "Items", actual: s.Plural()},
			"UnexportedVarName":       {expected: "item", actual: s.UnexportedVarName()},
			"PluralUnexportedVarName": {expected: "items", actual: s.PluralUnexportedVarName()},
			"RouteName":               {expected: "item", actual: s.RouteName()},
			"PluralRouteName":         {expected: "items", actual: s.PluralRouteName()},
		}

		for testName, zest := range expectationsMap {
			t.Run(testName, func(_t *testing.T) {
				assert.Equal(_t, zest.expected, zest.actual, "expected s.%s() to equal %q, not %q", testName, zest.expected, zest.actual)
			})
		}
	})

	T.Run("JournalEntry", func(t *testing.T) {
		s := FromSingularPascalCase("JournalEntry")

		expectationsMap := map[string]test{
			"Singular":                {expected: "JournalEntry", actual: s.Singular()},
			"Plural":                  {expected: "JournalEntries", actual: s.Plural()},
			"UnexportedVarName":       {expected: "journalEntry", actual: s.UnexportedVarName()},
			"PluralUnexportedVarName": {expected: "journalEntries", actual: s.PluralUnexportedVarName()},
			"RouteName":               {expected: "journal_entry", actual: s.RouteName()},
			"PluralRouteName":         {expected: "journal_entries", actual: s.PluralRouteName()},
		}

		for testName, zest := range expectationsMap {
			t.Run(testName, func(_t *testing.T) {
				assert.Equal(_t, zest.expected, zest.actual, "expected s.%s() to equal %q, not %q", testName, zest.expected, zest.actual)
			})
		}
	})
}
