package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

// buildOwnershipChain takes a series of names and returns a slice of datatypes with ownership between them.
// So for instance, if you provided `Forum`, `Subforum`, and `Post` as input, the output would be:
// 		[]DataType{
//			{
//				Name: wordsmith.FromSingularPascalCase("Forum"),
//			},
//			{
//				Name:            wordsmith.FromSingularPascalCase("Subforum"),
//				BelongsToEnumeration: wordsmith.FromSingularPascalCase("Forum"),
//			},
//			{
//				Name:            wordsmith.FromSingularPascalCase("Post"),
//				BelongsToEnumeration: wordsmith.FromSingularPascalCase("Subforum"),
//			},
//		}
func buildOwnershipChain(names ...string) (out []models.DataType) {
	for i, name := range names {
		if i == 0 {
			out = append(out,
				models.DataType{
					Name: wordsmith.FromSingularPascalCase(name),
				},
			)
		} else {
			out = append(out,
				models.DataType{
					Name:            wordsmith.FromSingularPascalCase(name),
					BelongsToStruct: wordsmith.FromSingularPascalCase(names[i-1]),
				},
			)
		}
	}

	return
}

func Test_wireParamFetchersTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireParamFetchersTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"testing"
)

func TestProvideItemsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemsServiceUserIDFetcher()
	})
}

func TestProvideItemsServiceItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemsServiceItemIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUsersServiceUserIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceUserIDFetcher()
	})
}

func TestProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceWebhookIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideOAuth2ClientsServiceClientIDFetcher(noop.ProvideNoopLogger())
	})
}

func Test_userIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		expected := exampleUser.ToSessionInfo()

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, expected),
		)

		actual := userIDFetcherFromRequestContext(req)
		assert.Equal(t, expected.UserID, actual)
	})

	T.Run("without attached value", func(t *testing.T) {
		req := buildRequest(t)
		actual := userIDFetcherFromRequestContext(req)

		assert.Zero(t, actual)
	})
}

func Test_buildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{items.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{items.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := wireParamFetchersTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	anotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/anotherthings"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	things "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/things"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	yetanotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/yetanotherthings"
	"testing"
)

func TestProvideThingsServiceThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideThingsServiceThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideAnotherThingsServiceThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideAnotherThingsServiceThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideAnotherThingsServiceAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideAnotherThingsServiceAnotherThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideYetAnotherThingsServiceThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideYetAnotherThingsServiceThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideYetAnotherThingsServiceAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideYetAnotherThingsServiceAnotherThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideYetAnotherThingsServiceYetAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideYetAnotherThingsServiceYetAnotherThingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUsersServiceUserIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceUserIDFetcher()
	})
}

func TestProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceWebhookIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideOAuth2ClientsServiceClientIDFetcher(noop.ProvideNoopLogger())
	})
}

func Test_userIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		expected := exampleUser.ToSessionInfo()

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, expected),
		)

		actual := userIDFetcherFromRequestContext(req)
		assert.Equal(t, expected.UserID, actual)
	})

	T.Run("without attached value", func(t *testing.T) {
		req := buildRequest(t)
		actual := userIDFetcherFromRequestContext(req)

		assert.Zero(t, actual)
	})
}

func Test_buildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{things.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{things.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamAnotherThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{anotherthings.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamAnotherThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{anotherthings.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamYetAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamYetAnotherThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{yetanotherthings.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamYetAnotherThingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{yetanotherthings.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideSomethingServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestProvideSomethingServiceUserIDFetcher(typ)

		expected := `
package example

import (
	"testing"
)

func TestProvideItemsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemsServiceUserIDFetcher()
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideSomethingServiceSomethingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestProvideSomethingServiceSomethingIDFetcher(typ)

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideItemsServiceItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemsServiceItemIDFetcher(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideSomethingServiceOwnerTypeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = buildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTestProvideSomethingServiceOwnerTypeIDFetcher(proj.LastDataType(), proj.DataTypes[1])

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideYetAnotherThingsServiceAnotherThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideYetAnotherThingsServiceAnotherThingIDFetcher(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideUsersServiceUserIDFetcher()

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUsersServiceUserIDFetcher(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideWebhooksServiceUserIDFetcher()

		expected := `
package example

import (
	"testing"
)

func TestProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceUserIDFetcher()
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideWebhooksServiceWebhookIDFetcher()

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceWebhookIDFetcher(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideOAuth2ClientsServiceClientIDFetcher()

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideOAuth2ClientsServiceClientIDFetcher(noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_userIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_userIDFetcherFromRequestContext(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func Test_userIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		expected := exampleUser.ToSessionInfo()

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, expected),
		)

		actual := userIDFetcherFromRequestContext(req)
		assert.Equal(t, expected.UserID, actual)
	})

	T.Run("without attached value", func(t *testing.T) {
		req := buildRequest(t)
		actual := userIDFetcherFromRequestContext(req)

		assert.Zero(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_buildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_buildRouteParamUserIDFetcher(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	"testing"
)

func Test_buildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{users.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_buildRouteParamSomethingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTest_buildRouteParamSomethingIDFetcher(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	"testing"
)

func Test_buildRouteParamItemIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{items.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{items.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_buildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_buildRouteParamWebhookIDFetcher(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"testing"
)

func Test_buildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooks.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_buildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_buildRouteParamOAuth2ClientIDFetcher(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	chi "github.com/go-chi/chi"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	"testing"
)

func Test_buildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clients.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
