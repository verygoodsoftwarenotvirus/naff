package httpserver

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireParamFetchersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireParamFetchersDotGo(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
	"strconv"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideItemsServiceUserIDFetcher,
		ProvideItemsServiceItemIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
	)
)

// ProvideItemsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideItemsServiceUserIDFetcher() items.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideItemsServiceItemIDFetcher provides an ItemIDFetcher.
func ProvideItemsServiceItemIDFetcher(logger v1.Logger) items.ItemIDFetcher {
	return buildRouteParamItemIDFetcher(logger)
}

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooks.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger v1.Logger) webhooks.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger v1.Logger) oauth2clients.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// userIDFetcherFromRequestContext fetches a user ID from a request routed by chi.
func userIDFetcherFromRequestContext(req *http.Request) uint64 {
	if si, ok := req.Context().Value(v11.SessionInfoKey).(*v11.SessionInfo); ok && si != nil {
		return si.UserID
	}
	return 0
}

// buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi.
func buildRouteParamUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return func(req *http.Request) uint64 {
		u, err := strconv.ParseUint(chi.URLParam(req, users.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching user ID from request")
		}
		return u
	}
}

// buildRouteParamItemIDFetcher builds a function that fetches a ItemID from a request routed by chi.
func buildRouteParamItemIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, items.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ItemID from request")
		}
		return u
	}
}

// buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi.
func buildRouteParamWebhookIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, webhooks.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching WebhookID from request")
		}
		return u
	}
}

// buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.
func buildRouteParamOAuth2ClientIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, oauth2clients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching OAuth2ClientID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := wireParamFetchersDotGo(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	anotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/anotherthings"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	things "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/things"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	yetanotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/yetanotherthings"
	"net/http"
	"strconv"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideThingsServiceThingIDFetcher,
		ProvideAnotherThingsServiceThingIDFetcher,
		ProvideAnotherThingsServiceAnotherThingIDFetcher,
		ProvideYetAnotherThingsServiceThingIDFetcher,
		ProvideYetAnotherThingsServiceAnotherThingIDFetcher,
		ProvideYetAnotherThingsServiceYetAnotherThingIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
	)
)

// ProvideThingsServiceThingIDFetcher provides a ThingIDFetcher.
func ProvideThingsServiceThingIDFetcher(logger v1.Logger) things.ThingIDFetcher {
	return buildRouteParamThingIDFetcher(logger)
}

// ProvideThingsServiceThingIDFetcher provides a ThingIDFetcher.
func ProvideThingsServiceThingIDFetcher(logger v1.Logger) things.ThingIDFetcher {
	return buildRouteParamThingIDFetcher(logger)
}

// ProvideAnotherThingsServiceAnotherThingIDFetcher provides an AnotherThingIDFetcher.
func ProvideAnotherThingsServiceAnotherThingIDFetcher(logger v1.Logger) anotherthings.AnotherThingIDFetcher {
	return buildRouteParamAnotherThingIDFetcher(logger)
}

// ProvideThingsServiceThingIDFetcher provides a ThingIDFetcher.
func ProvideThingsServiceThingIDFetcher(logger v1.Logger) things.ThingIDFetcher {
	return buildRouteParamThingIDFetcher(logger)
}

// ProvideAnotherThingsServiceAnotherThingIDFetcher provides an AnotherThingIDFetcher.
func ProvideAnotherThingsServiceAnotherThingIDFetcher(logger v1.Logger) anotherthings.AnotherThingIDFetcher {
	return buildRouteParamAnotherThingIDFetcher(logger)
}

// ProvideYetAnotherThingsServiceYetAnotherThingIDFetcher provides a YetAnotherThingIDFetcher.
func ProvideYetAnotherThingsServiceYetAnotherThingIDFetcher(logger v1.Logger) yetanotherthings.YetAnotherThingIDFetcher {
	return buildRouteParamYetAnotherThingIDFetcher(logger)
}

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooks.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger v1.Logger) webhooks.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger v1.Logger) oauth2clients.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// userIDFetcherFromRequestContext fetches a user ID from a request routed by chi.
func userIDFetcherFromRequestContext(req *http.Request) uint64 {
	if si, ok := req.Context().Value(v11.SessionInfoKey).(*v11.SessionInfo); ok && si != nil {
		return si.UserID
	}
	return 0
}

// buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi.
func buildRouteParamUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return func(req *http.Request) uint64 {
		u, err := strconv.ParseUint(chi.URLParam(req, users.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching user ID from request")
		}
		return u
	}
}

// buildRouteParamThingIDFetcher builds a function that fetches a ThingID from a request routed by chi.
func buildRouteParamThingIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, things.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ThingID from request")
		}
		return u
	}
}

// buildRouteParamAnotherThingIDFetcher builds a function that fetches a AnotherThingID from a request routed by chi.
func buildRouteParamAnotherThingIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, anotherthings.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching AnotherThingID from request")
		}
		return u
	}
}

// buildRouteParamYetAnotherThingIDFetcher builds a function that fetches a YetAnotherThingID from a request routed by chi.
func buildRouteParamYetAnotherThingIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, yetanotherthings.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching YetAnotherThingID from request")
		}
		return u
	}
}

// buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi.
func buildRouteParamWebhookIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, webhooks.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching WebhookID from request")
		}
		return u
	}
}

// buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.
func buildRouteParamOAuth2ClientIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, oauth2clients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching OAuth2ClientID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireParamFetchersVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWireParamFetchersVarDeclarations(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideItemsServiceUserIDFetcher,
		ProvideItemsServiceItemIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildWireParamFetchersVarDeclarations(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideThingsServiceThingIDFetcher,
		ProvideAnotherThingsServiceThingIDFetcher,
		ProvideAnotherThingsServiceAnotherThingIDFetcher,
		ProvideYetAnotherThingsServiceThingIDFetcher,
		ProvideYetAnotherThingsServiceAnotherThingIDFetcher,
		ProvideYetAnotherThingsServiceYetAnotherThingIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideSomethingServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildProvideSomethingServiceUserIDFetcher(proj, typ)

		expected := `
package example

import (
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
)

// ProvideItemsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideItemsServiceUserIDFetcher() items.UserIDFetcher {
	return userIDFetcherFromRequestContext
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideSomethingServiceThingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildProvideSomethingServiceThingIDFetcher(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
)

// ProvideItemsServiceItemIDFetcher provides an ItemIDFetcher.
func ProvideItemsServiceItemIDFetcher(logger v1.Logger) items.ItemIDFetcher {
	return buildRouteParamItemIDFetcher(logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideSomethingServiceOwnerTypeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildProvideSomethingServiceOwnerTypeIDFetcher(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
)

// ProvideItemsServiceItemIDFetcher provides an ItemIDFetcher.
func ProvideItemsServiceItemIDFetcher(logger v1.Logger) items.ItemIDFetcher {
	return buildRouteParamItemIDFetcher(logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideUsersServiceUserIDFetcher(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
)

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideWebhooksServiceUserIDFetcher(proj)

		expected := `
package example

import (
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
)

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooks.UserIDFetcher {
	return userIDFetcherFromRequestContext
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideWebhooksServiceWebhookIDFetcher(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
)

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger v1.Logger) webhooks.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideOAuth2ClientsServiceClientIDFetcher(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
)

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger v1.Logger) oauth2clients.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUserIDFetcherFromRequestContext(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// userIDFetcherFromRequestContext fetches a user ID from a request routed by chi.
func userIDFetcherFromRequestContext(req *http.Request) uint64 {
	if si, ok := req.Context().Value(v1.SessionInfoKey).(*v1.SessionInfo); ok && si != nil {
		return si.UserID
	}
	return 0
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildRouteParamUserIDFetcher(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	"net/http"
	"strconv"
)

// buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi.
func buildRouteParamUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return func(req *http.Request) uint64 {
		u, err := strconv.ParseUint(chi.URLParam(req, users.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching user ID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRouteParamSomethingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildRouteParamSomethingIDFetcher(proj, typ)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	"net/http"
	"strconv"
)

// buildRouteParamItemIDFetcher builds a function that fetches a ItemID from a request routed by chi.
func buildRouteParamItemIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, items.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ItemID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildRouteParamWebhookIDFetcher(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
	"strconv"
)

// buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi.
func buildRouteParamWebhookIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, webhooks.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching WebhookID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildRouteParamOAuth2ClientIDFetcher(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	"net/http"
	"strconv"
)

// buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.
func buildRouteParamOAuth2ClientIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, oauth2clients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching OAuth2ClientID from request")
		}
		return u
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
