package models

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func Test_buildFakeVarName(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		exampleInput := "fart"

		expected := "exampleFart"
		actual := buildFakeVarName(exampleInput)

		assert.Equal(t, expected, actual)
	})
}

func Test_ctxParam(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(
			jen.Func().ID("test").Params(
				ctxParam(),
			).Body(),
		)

		var b bytes.Buffer
		assert.NoError(t, out.Render(&b))

		expected := `
package main

import (
	"context"
)

func test(ctx context.Context) {}
`
		actual := "\n" + b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_ctxVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(
			jen.Func().ID("test").Params().Body(
				ctxVar(),
			),
		)

		var b bytes.Buffer
		assert.NoError(t, out.Render(&b))

		expected := `
package main

import ()

func test() {
	ctx
}
`
		actual := "\n" + b.String()

		assert.Equal(t, expected, actual)
	})
}

// DataType tests

func TestDataType_OwnedByAUserAtSomeLevel(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToAccount: true,
		}
		p := &Project{
			DataTypes: []DataType{
				dt,
			},
		}

		assert.True(t, dt.OwnedByAUserAtSomeLevel(p))
	})

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			DataTypes: BuildOwnershipChain("A", "B", "C"),
		}
		p.DataTypes[0].BelongsToAccount = true

		assert.True(t, p.LastDataType().OwnedByAUserAtSomeLevel(p))
	})
}

func TestDataType_RestrictedToUserAtSomeLevel(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToAccount:           true,
			RestrictedToAccountMembers: true,
		}
		p := &Project{
			DataTypes: []DataType{
				dt,
			},
		}

		assert.True(t, dt.RestrictedToAccountAtSomeLevel(p))
	})

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		dtA := DataType{
			Name:                       wordsmith.FromSingularPascalCase("A"),
			BelongsToAccount:           true,
			RestrictedToAccountMembers: true,
		}
		dtB := DataType{
			Name:            wordsmith.FromSingularPascalCase("B"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
		}
		dtC := DataType{
			Name:            wordsmith.FromSingularPascalCase("C"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
		}
		p := &Project{
			DataTypes: []DataType{
				dtA,
				dtB,
				dtC,
			},
		}

		assert.True(t, dtC.RestrictedToAccountAtSomeLevel(p))
	})
}

func TestDataType_MultipleOwnersBelongingToUser(T *testing.T) {
	T.Parallel()

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		dtA := DataType{
			Name:             wordsmith.FromSingularPascalCase("A"),
			BelongsToAccount: true,
		}
		dtB := DataType{
			Name:             wordsmith.FromSingularPascalCase("B"),
			BelongsToAccount: true,
			BelongsToStruct:  wordsmith.FromSingularPascalCase("A"),
		}
		dtC := DataType{
			Name:             wordsmith.FromSingularPascalCase("C"),
			BelongsToAccount: true,
			BelongsToStruct:  wordsmith.FromSingularPascalCase("B"),
		}
		p := &Project{
			DataTypes: []DataType{
				dtA,
				dtB,
				dtC,
			},
		}

		assert.True(t, dtC.MultipleOwnersBelongingToUser(p))
	})
}

func TestDataType_buildGetSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		p.DataTypes[0].BelongsToAccount = true
		p.DataTypes[0].RestrictedToAccountMembers = true

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID, userID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().buildGetSomethingParams(p, true))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildArchiveSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:            wordsmith.FromSingularPascalCase("Thing"),
			Name:                       wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount:           true,
			RestrictedToAccountMembers: true,
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, userID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.buildArchiveSomethingParams())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionArchiveMethodParams())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientArchiveMethodParams())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierArchiveMethodParams())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func example(thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierArchiveQueryMethodParams())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func example(thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierRetrievalQueryMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func example(thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierExistenceQueryMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_ModifyQueryBuildingStatementWithJoinClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing", "EvenStillAnotherThing")

		result := p.LastDataType().ModifyQueryBuildingStatementWithJoinClauses(p, jen.ID("something"))

		expected := `
package main

import ()

func main() {
	something.
		Join(yetAnotherThingsOnEvenStillAnotherThingsJoinClause).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause)
}
`
		actual := renderIndependentStatementToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildJoinClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name:            wordsmith.FromSingularPascalCase("Thing"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("SomethingElse"),
		}

		expected := "table1 ON table2.belongs_to_table3=table1.id"
		actual := dt.buildJoinClause("table1", "table2", "table3")

		assert.Equal(t, expected, actual)
	})

	T.Run("panics on non-ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur!")
			}
		}()

		dt.buildJoinClause("table1", "table2", "table3")
	})
}

func TestDataType_ModifyQueryBuilderWithJoinClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		s := squirrel.Select("*").From("fart")

		result := p.LastDataType().ModifyQueryBuilderWithJoinClauses(p, s)

		expected := "SELECT * FROM fart JOIN another_things ON yet_another_things.belongs_to_another_thing=another_things.id JOIN things ON another_things.belongs_to_thing=things.id"
		actual, _, err := result.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDBQuerierSingleInstanceQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		result := dt.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p)

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", thingsTableName, idColumn): thingID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, result)

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		result := p.LastDataType().buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p)

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
		fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
		fmt.Sprintf("%s.%s", thingsTableName, thingsUserOwnershipColumn):                      userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsUserOwnershipColumn):        userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsUserOwnershipColumn):  userID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		result := dt.BuildDBQuerierExistenceQueryMethodConditionalClauses(p)

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", thingsTableName, idColumn): thingID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		results := dt.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE things.id = ?"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		results := p.LastDataType().buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE another_things.belongs_to_thing = ? AND another_things.belongs_to_user = ? AND another_things.id = ? AND things.belongs_to_user = ? AND things.id = ? AND yet_another_things.belongs_to_another_thing = ? AND yet_another_things.belongs_to_user = ? AND yet_another_things.id = ?"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		results := dt.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE things.id = ?"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		results := dt.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE things.id = ?"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		results := dt.BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE things.archived_on IS NULL"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		results := p.LastDataType().BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(p)
		qb := squirrel.Select("*").From("farts").Where(results)

		expected := "SELECT * FROM farts WHERE another_things.belongs_to_thing = ? AND another_things.belongs_to_user = ? AND another_things.id = ? AND things.belongs_to_user = ? AND things.id = ? AND yet_another_things.archived_on IS NULL AND yet_another_things.belongs_to_another_thing = ? AND yet_another_things.belongs_to_user = ?"
		actual, _, err := qb.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", thingsTableName, idColumn): thingID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, dt.BuildDBQuerierRetrievalQueryMethodConditionalClauses(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
		fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
		fmt.Sprintf("%s.%s", thingsTableName, thingsUserOwnershipColumn):                      userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsUserOwnershipColumn):        userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsUserOwnershipColumn):  userID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, p.LastDataType().BuildDBQuerierRetrievalQueryMethodConditionalClauses(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", thingsTableName, archivedOnColumn): nil,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, dt.BuildDBQuerierListRetrievalQueryMethodConditionalClauses(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleMap := map[string]interface{}{
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn):                     nil,
		fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
		fmt.Sprintf("%s.%s", thingsTableName, thingsUserOwnershipColumn):                      userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsUserOwnershipColumn):        userID,
		fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsUserOwnershipColumn):  userID,
	}
}
`
		actual := renderMapEntriesWithStringKeysToString(t, p.LastDataType().BuildDBQuerierListRetrievalQueryMethodConditionalClauses(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.buildGetSomethingArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownerships", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, yetAnotherThingID, userID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildGetSomethingArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildArchiveSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.buildArchiveSomethingArgs())

		assert.Equal(t, expected, actual)
	})

	T.Run("with ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thing"),
			Name:             wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount: true,
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, userID)
}
`
		actual := renderCallArgsToString(t, dt.buildArchiveSomethingArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientExistenceMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientExistenceMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientRetrievalMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientArchiveMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientArchiveMethodCallArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierExistenceQueryBuildingArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierRetrievalQueryBuildingArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierArchiveQueryBuildingArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionExistenceMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildInterfaceDefinitionExistenceMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildInterfaceDefinitionRetrievalMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionArchiveMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildInterfaceDefinitionArchiveMethodCallArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetSomethingArgsWithExampleVariables(T *testing.T) {
	T.Parallel()

	T.Run("while including context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.buildGetSomethingArgsWithExampleVariables(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without including context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.buildGetSomethingArgsWithExampleVariables(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildGetSomethingArgsWithExampleVariables(p, true))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientRetrievalTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildHTTPClientRetrievalTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildSingleInstanceQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.buildSingleInstanceQueryTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleYetAnotherThing.BelongsToAccount)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildSingleInstanceQueryTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForServiceRouteExistenceCheck(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForServiceRouteExistenceCheck(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, yetAnotherThingID, userID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForServiceRouteExistenceCheck(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with forums", func(t *testing.T) {
		t.Parallel()

		p := buildExampleForumsListProject()
		dt := p.DataTypes[3]

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleForum.ID, exampleSubforum.ID, exampleThread.ID, examplePost.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierRetrievalQueryTestCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	filter := fake.BuildFleshedOutQueryFilter()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	filter := fake.BuildFleshedOutQueryFilter()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildGetSomethingLogValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	map[string]interface{}{
		"thing_id": thingID,
	}
}
`
		actual := renderIndependentStatementToString(t, dt.BuildGetSomethingLogValues(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	map[string]interface{}{
		"thing_id":             thingID,
		"another_thing_id":     anotherThingID,
		"yet_another_thing_id": yetAnotherThingID,
		"user_id":              userID,
	}
}
`
		actual := renderIndependentStatementToString(t, p.LastDataType().BuildGetSomethingLogValues(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildGetListOfSomethingLogValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	map[string]interface{}{
		"thing_id":         thingID,
		"another_thing_id": anotherThingID,
	}
}
`
		actual := renderIndependentStatementToString(t, p.LastDataType().BuildGetListOfSomethingLogValues(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	map[string]interface{}{
		"thing_id":         thingID,
		"another_thing_id": anotherThingID,
		"user_id":          userID,
	}
}
`
		actual := renderIndependentStatementToString(t, p.LastDataType().BuildGetListOfSomethingLogValues(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetListOfSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, filter *QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.buildGetListOfSomethingParams(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.buildGetListOfSomethingParams(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, userID uint64, filter *QueryFilter) {}
`

		actual := renderFunctionParamsToString(t, p.LastDataType().buildGetListOfSomethingParams(p, true))
		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildMockDataManagerListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, filter *QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierListRetrievalQueryBuildingMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildCreateSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, input *ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.buildCreateSomethingParams(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.buildCreateSomethingParams(p, false))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockInterfaceDefinitionCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildMockInterfaceDefinitionCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, input *ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, input *Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierCreationQueryBuildingMethodParams(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierCreationQueryBuildingMethodParams(p, false))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildCreateSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, input)
}
`
		actual := renderCallArgsToString(t, dt.buildCreateSomethingArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockInterfaceDefinitionCreationMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, input)
}
`
		actual := renderCallArgsToString(t, dt.BuildMockInterfaceDefinitionCreationMethodCallArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(input)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierCreationMethodQueryBuildingArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfUpdateQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierTestOfUpdateQueryBuilder())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfArchiveQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierTestOfArchiveQueryBuilder())

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thing"),
			Name:             wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount: true,
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierTestOfArchiveQueryBuilder())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfUpdateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForDBQuerierTestOfUpdateMethod())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodArgsToUseFromMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleInput)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierCreationMethodArgsToUseFromMethodTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsToUseForDBQuerierCreationQueryBuildingTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsToUseForDBQuerierCreationQueryBuildingTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientCreationMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, input)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientCreationMethodCallArgs())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, updated *Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.buildUpdateSomethingParams(p, "updated", true))

		assert.Equal(t, expected, actual)
	})

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, updated *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.buildUpdateSomethingParams(p, "updated", false))

		assert.Equal(t, expected, actual)
	})

	T.Run("panics with empty input", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		renderFunctionParamsToString(t, p.LastDataType().buildUpdateSomethingParams(p, "", true))
	})

}

func TestDataType_BuildDBClientUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, updated *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientUpdateMethodParams(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, updated *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierUpdateMethodParams(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(updated *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierUpdateQueryBuildingMethodParams(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, updated *Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionUpdateMethodParams(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, updated *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildMockDataManagerUpdateMethodParams(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingArgsWithExampleVars(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, updated)
}
`
		actual := renderCallArgsToString(t, dt.buildUpdateSomethingArgsWithExampleVars(p, "updated"))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, Thing.ID, updated)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildUpdateSomethingArgsWithExampleVars(p, "updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, updated)
}
`
		actual := renderCallArgsToString(t, dt.buildUpdateSomethingArgs("updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientUpdateMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, updated)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientUpdateMethodCallArgs("updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateMethodArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(updated)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierUpdateMethodArgs("updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerUpdateMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, updated)
}
`
		actual := renderCallArgsToString(t, dt.BuildMockDataManagerUpdateMethodCallArgs("updated"))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetListOfSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.buildGetListOfSomethingArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, userID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildGetListOfSomethingArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientListRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBClientListRetrievalMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalMethodArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildDBQuerierListRetrievalMethodArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerListRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildMockDataManagerListRetrievalMethodCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWithOwnerStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildVarDeclarationsOfDependentStructsWithOwnerStruct(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildVarDeclarationsOfDependentStructsWithOwnerStruct(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForDBQueriersExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForDBQueriersExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForDBQueriersCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForDBQueriersCreationMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientRetrievalMethodTestDependentObjects(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildHTTPClientRetrievalMethodTestDependentObjects(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildDependentObjectsForHTTPClientArchiveMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDependentObjectsForHTTPClientListRetrievalTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildDependentObjectsForHTTPClientListRetrievalTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientListRetrievalTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDependentObjectsForHTTPClientListRetrievalTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsForUpdateFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildVarDeclarationsOfDependentStructsForUpdateFunction(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDependentObjectsForHTTPClientUpdateMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildDependentObjectsForHTTPClientCreationMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things/%d"
		actual := p.LastDataType().BuildFormatStringForHTTPClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things/%d"
		actual := p.LastDataType().BuildFormatStringForHTTPClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things/%d"
		actual := p.LastDataType().BuildFormatStringForHTTPClientUpdateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things/%d"
		actual := p.LastDataType().BuildFormatStringForHTTPClientArchiveMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientListMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things"
		actual := p.LastDataType().BuildFormatStringForHTTPClientListMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientSearchMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/yet_another_things/search"
		actual := p.LastDataType().BuildFormatStringForHTTPClientSearchMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientCreateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := "/api/v1/things/%d/another_things/%d/yet_another_things"
		actual := p.LastDataType().BuildFormatStringForHTTPClientCreateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildFormatCallArgsForHTTPClientRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildFormatCallArgsForHTTPClientRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildFormatCallArgsForHTTPClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildFormatCallArgsForHTTPClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientListMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		actual := renderCallArgsToString(t, dt.BuildFormatCallArgsForHTTPClientListMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildFormatCallArgsForHTTPClientListMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		actual := renderCallArgsToString(t, dt.BuildFormatCallArgsForHTTPClientCreationMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildFormatCallArgsForHTTPClientCreationMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientUpdateTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildFormatCallArgsForHTTPClientUpdateTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleYetAnotherThing.BelongsToAnotherThing, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildFormatCallArgsForHTTPClientUpdateTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientExistenceRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, yetAnotherThingID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientExistenceRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientExistenceRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientExistenceRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientExistenceRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientExistenceMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientExistenceMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientExistenceMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientCreateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, input)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientCreateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, input)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientCreateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientRetrievalRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientRetrievalRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, yetAnotherThingID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientRetrievalRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientRetrievalRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientRetrievalRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientRetrievalRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientRetrievalMethod(T *testing.T) {
	T.Parallel()

	T.Run("as call", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func example(ctx, thingID) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientRetrievalMethod(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("not as call", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientRetrievalMethod(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("as call with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func example(ctx, thingID, anotherThingID, yetAnotherThingID) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientRetrievalMethod(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("not as call with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientRetrievalMethod(p, false))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientCreateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientCreateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID uint64, input *types.YetAnotherThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientCreateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientCreateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, input *types.ThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientCreateMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID uint64, input *types.YetAnotherThingCreationInput) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientCreateMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientUpdateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thing *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientUpdateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID uint64, yetAnotherThing *types.YetAnotherThing) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientUpdateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientUpdateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thing)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientUpdateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, yetAnotherThing)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientUpdateRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientUpdateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thing *types.Thing) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientUpdateMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID uint64, yetAnotherThing *types.YetAnotherThing) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientUpdateMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientArchiveRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientArchiveRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientArchiveRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientArchiveRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, yetAnotherThingID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientArchiveRequestBuildingMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientArchiveMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientArchiveMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientArchiveMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildParamsForMethodThatHandlesAnInstanceWithStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func example(ctx, exampleThing.ID) {}
`
		actual := renderFunctionParamsToString(t, dt.buildParamsForMethodThatHandlesAnInstanceWithStructs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func example(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().buildParamsForMethodThatHandlesAnInstanceWithStructs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientArchiveMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientCreationRequestBuildingMethodArgsForTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleInput)
}
`
		actual := renderCallArgsToString(t, dt.BuildHTTPClientCreationRequestBuildingMethodArgsForTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleInput)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildHTTPClientCreationRequestBuildingMethodArgsForTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientCreationMethodArgsForTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleInput)
}
`
		actual := renderCallArgsToString(t, dt.BuildHTTPClientCreationMethodArgsForTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleInput)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildHTTPClientCreationMethodArgsForTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientListRequestMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildArgsForHTTPClientListRequestMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID, anotherThingID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildArgsForHTTPClientListRequestMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientListRequestMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientListRequestMethod(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID, anotherThingID uint64, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientListRequestMethod(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientMethodThatFetchesAList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildParamsForHTTPClientMethodThatFetchesAList(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	types "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types"
)

func example(ctx context.Context, thingID, anotherThingID uint64, filter *types.QueryFilter) {}
`
		actual := renderFunctionParamsToString(t, p.LastDataType().BuildParamsForHTTPClientMethodThatFetchesAList(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForHTTPClientListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForHTTPClientListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleYetAnotherThing)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForHTTPClientUpdateMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleYetAnotherThing)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForHTTPClientUpdateMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecs(T *testing.T) {
	T.Parallel()

	T.Run("creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildRequisiteFakeVarDecs(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildRequisiteFakeVarDecs(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleUser := fake.BuildFakeUser()
	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildRequisiteFakeVarDecs(p, true))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecForModifierFuncs(T *testing.T) {
	T.Parallel()

	T.Run("creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildRequisiteFakeVarDecForModifierFuncs(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.buildRequisiteFakeVarDecForModifierFuncs(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleUser := fake.BuildFakeUser()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildRequisiteFakeVarDecForModifierFuncs(p, true))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildRequisiteFakeVarsForDBClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientCreateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildRequisiteFakeVarsForDBClientCreateMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	exampleUser := fake.BuildFakeUser()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBClientCreateMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	var expected error

	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildRequisiteFakeVarsForDBClientArchiveMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()

	var expected error

	exampleUser := fake.BuildFakeUser()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBClientArchiveMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
}
`
		actual := renderVariableDeclarationsToString(t, dt.BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecsForListFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {}
`
		actual := renderVariableDeclarationsToString(t, dt.buildRequisiteFakeVarDecsForListFunction(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			if i != len(p.DataTypes)-1 {
				p.DataTypes[i].BelongsToAccount = true
				p.DataTypes[i].RestrictedToAccountMembers = true
			}
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildRequisiteFakeVarDecsForListFunction(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("with filter", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	filter := fake.BuildFleshedOutQueryFilter()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p, true))

		assert.Equal(t, expected, actual)
	})

	T.Run("without filter", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p, false))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	filter := fake.BuildFleshedOutQueryFilter()
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p, true))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarCallArgsForCreation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleThing.ID
	exampleAnotherThing.ID
	exampleYetAnotherThing.ID
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgsForCreation(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleThing.ID
	exampleAnotherThing.ID
	exampleYetAnotherThing.ID
	exampleYetAnotherThing.BelongsToAccount
}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgsForCreation(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("type restricted to user", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleYetAnotherThing.BelongsToAccount)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgs(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("belonging to type restricted to user", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		p.DataTypes[0].BelongsToAccount = true
		p.DataTypes[0].RestrictedToAccountMembers = true

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgs(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with user restriction", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		p.DataTypes[0].BelongsToAccount = true
		p.DataTypes[0].RestrictedToAccountMembers = true

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest())

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thing"),
			Name:             wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount: true,
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest())

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thing"),
			Name:             wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount: true,
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleAnotherThing.BelongsToThing, exampleAnotherThing.ID, exampleAnotherThing.BelongsToAccount)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with user restriction", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		p.DataTypes[2].BelongsToAccount = true
		p.DataTypes[2].RestrictedToAccountMembers = true

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, filter)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID, filter)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest())

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thing"),
			Name:             wordsmith.FromSingularPascalCase("AnotherThing"),
			BelongsToAccount: true,
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, dt.BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, exampleInput)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForDBClientCreationMethodTest())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForDBClientListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with user restriction", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		p.DataTypes[0].RestrictedToAccountMembers = true
		p.DataTypes[0].BelongsToAccount = true

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing.ID, exampleAnotherThing.ID, exampleUser.ID)
}
`
		actual := renderCallArgsToString(t, p.LastDataType().BuildCallArgsForDBClientListRetrievalMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteVarsForDBClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()
	var expected error

	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()

}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteVarsForDBClientUpdateMethodTest(p))

		assert.Equal(t, expected, actual)
	})

	T.Run("with multiple ownership", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes = BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range p.DataTypes {
			p.DataTypes[i].BelongsToAccount = true
			p.DataTypes[i].RestrictedToAccountMembers = true
		}

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
)

func main() {
	ctx := context.Background()
	var expected error

	exampleUser := fake.BuildFakeUser()
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAccount = exampleUser.ID

}
`
		actual := renderVariableDeclarationsToString(t, p.LastDataType().BuildRequisiteVarsForDBClientUpdateMethodTest(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(exampleThing)
}
`
		actual := renderCallArgsToString(t, dt.BuildCallArgsForDBClientUpdateMethodTest())

		assert.Equal(t, expected, actual)
	})
}
