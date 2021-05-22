package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
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
		ProvideWebhooksServiceWebhookIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideItemsServiceItemIDFetcher,
		ProvideItemsServiceUserIDFetcher,
	)
)

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger v1.Logger) oauth2clients.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger v1.Logger) webhooks.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooks.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideItemsServiceItemIDFetcher provides an ItemIDFetcher.
func ProvideItemsServiceItemIDFetcher(logger v1.Logger) items.ItemIDFetcher {
	return buildRouteParamItemIDFetcher(logger)
}

// ProvideItemsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideItemsServiceUserIDFetcher() items.UserIDFetcher {
	return userIDFetcherFromRequestContext
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

	T.Run("with forums app", func(t *testing.T) {
		proj := testprojects.BuildForumsApp()
		x := wireParamFetchersDotGo(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	forums "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/forums"
	notifications "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/notifications"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	postreactions "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/postreactions"
	posts "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/posts"
	reactionicons "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/reactionicons"
	subforums "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/subforums"
	threads "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/threads"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
	"strconv"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideForumsServiceForumIDFetcher,
		ProvideSubforumsServiceForumIDFetcher,
		ProvideSubforumsServiceSubforumIDFetcher,
		ProvideThreadsServiceForumIDFetcher,
		ProvideThreadsServiceSubforumIDFetcher,
		ProvideThreadsServiceThreadIDFetcher,
		ProvideThreadsServiceUserIDFetcher,
		ProvidePostsServiceForumIDFetcher,
		ProvidePostsServiceSubforumIDFetcher,
		ProvidePostsServiceThreadIDFetcher,
		ProvidePostsServicePostIDFetcher,
		ProvidePostsServiceUserIDFetcher,
		ProvideReactionIconsServiceReactionIconIDFetcher,
		ProvidePostReactionsServiceForumIDFetcher,
		ProvidePostReactionsServiceSubforumIDFetcher,
		ProvidePostReactionsServiceThreadIDFetcher,
		ProvidePostReactionsServicePostIDFetcher,
		ProvidePostReactionsServicePostReactionIDFetcher,
		ProvidePostReactionsServiceUserIDFetcher,
		ProvideNotificationsServiceNotificationIDFetcher,
		ProvideNotificationsServiceUserIDFetcher,
	)
)

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger v1.Logger) users.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger v1.Logger) oauth2clients.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger v1.Logger) webhooks.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooks.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideForumsServiceForumIDFetcher provides a ForumIDFetcher.
func ProvideForumsServiceForumIDFetcher(logger v1.Logger) forums.ForumIDFetcher {
	return buildRouteParamForumIDFetcher(logger)
}

// ProvideSubforumsServiceForumIDFetcher provides a ForumIDFetcher.
func ProvideSubforumsServiceForumIDFetcher(logger v1.Logger) subforums.ForumIDFetcher {
	return buildRouteParamForumIDFetcher(logger)
}

// ProvideSubforumsServiceSubforumIDFetcher provides a SubforumIDFetcher.
func ProvideSubforumsServiceSubforumIDFetcher(logger v1.Logger) subforums.SubforumIDFetcher {
	return buildRouteParamSubforumIDFetcher(logger)
}

// ProvideThreadsServiceForumIDFetcher provides a ForumIDFetcher.
func ProvideThreadsServiceForumIDFetcher(logger v1.Logger) threads.ForumIDFetcher {
	return buildRouteParamForumIDFetcher(logger)
}

// ProvideThreadsServiceSubforumIDFetcher provides a SubforumIDFetcher.
func ProvideThreadsServiceSubforumIDFetcher(logger v1.Logger) threads.SubforumIDFetcher {
	return buildRouteParamSubforumIDFetcher(logger)
}

// ProvideThreadsServiceThreadIDFetcher provides a ThreadIDFetcher.
func ProvideThreadsServiceThreadIDFetcher(logger v1.Logger) threads.ThreadIDFetcher {
	return buildRouteParamThreadIDFetcher(logger)
}

// ProvideThreadsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideThreadsServiceUserIDFetcher() threads.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvidePostsServiceForumIDFetcher provides a ForumIDFetcher.
func ProvidePostsServiceForumIDFetcher(logger v1.Logger) posts.ForumIDFetcher {
	return buildRouteParamForumIDFetcher(logger)
}

// ProvidePostsServiceSubforumIDFetcher provides a SubforumIDFetcher.
func ProvidePostsServiceSubforumIDFetcher(logger v1.Logger) posts.SubforumIDFetcher {
	return buildRouteParamSubforumIDFetcher(logger)
}

// ProvidePostsServiceThreadIDFetcher provides a ThreadIDFetcher.
func ProvidePostsServiceThreadIDFetcher(logger v1.Logger) posts.ThreadIDFetcher {
	return buildRouteParamThreadIDFetcher(logger)
}

// ProvidePostsServicePostIDFetcher provides a PostIDFetcher.
func ProvidePostsServicePostIDFetcher(logger v1.Logger) posts.PostIDFetcher {
	return buildRouteParamPostIDFetcher(logger)
}

// ProvidePostsServiceUserIDFetcher provides a UserIDFetcher.
func ProvidePostsServiceUserIDFetcher() posts.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideReactionIconsServiceReactionIconIDFetcher provides a ReactionIconIDFetcher.
func ProvideReactionIconsServiceReactionIconIDFetcher(logger v1.Logger) reactionicons.ReactionIconIDFetcher {
	return buildRouteParamReactionIconIDFetcher(logger)
}

// ProvidePostReactionsServiceForumIDFetcher provides a ForumIDFetcher.
func ProvidePostReactionsServiceForumIDFetcher(logger v1.Logger) postreactions.ForumIDFetcher {
	return buildRouteParamForumIDFetcher(logger)
}

// ProvidePostReactionsServiceSubforumIDFetcher provides a SubforumIDFetcher.
func ProvidePostReactionsServiceSubforumIDFetcher(logger v1.Logger) postreactions.SubforumIDFetcher {
	return buildRouteParamSubforumIDFetcher(logger)
}

// ProvidePostReactionsServiceThreadIDFetcher provides a ThreadIDFetcher.
func ProvidePostReactionsServiceThreadIDFetcher(logger v1.Logger) postreactions.ThreadIDFetcher {
	return buildRouteParamThreadIDFetcher(logger)
}

// ProvidePostReactionsServicePostIDFetcher provides a PostIDFetcher.
func ProvidePostReactionsServicePostIDFetcher(logger v1.Logger) postreactions.PostIDFetcher {
	return buildRouteParamPostIDFetcher(logger)
}

// ProvidePostReactionsServicePostReactionIDFetcher provides a PostReactionIDFetcher.
func ProvidePostReactionsServicePostReactionIDFetcher(logger v1.Logger) postreactions.PostReactionIDFetcher {
	return buildRouteParamPostReactionIDFetcher(logger)
}

// ProvidePostReactionsServiceUserIDFetcher provides a UserIDFetcher.
func ProvidePostReactionsServiceUserIDFetcher() postreactions.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideNotificationsServiceNotificationIDFetcher provides a NotificationIDFetcher.
func ProvideNotificationsServiceNotificationIDFetcher(logger v1.Logger) notifications.NotificationIDFetcher {
	return buildRouteParamNotificationIDFetcher(logger)
}

// ProvideNotificationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideNotificationsServiceUserIDFetcher() notifications.UserIDFetcher {
	return userIDFetcherFromRequestContext
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

// buildRouteParamForumIDFetcher builds a function that fetches a ForumID from a request routed by chi.
func buildRouteParamForumIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, forums.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ForumID from request")
		}
		return u
	}
}

// buildRouteParamSubforumIDFetcher builds a function that fetches a SubforumID from a request routed by chi.
func buildRouteParamSubforumIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, subforums.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching SubforumID from request")
		}
		return u
	}
}

// buildRouteParamThreadIDFetcher builds a function that fetches a ThreadID from a request routed by chi.
func buildRouteParamThreadIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, threads.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ThreadID from request")
		}
		return u
	}
}

// buildRouteParamPostIDFetcher builds a function that fetches a PostID from a request routed by chi.
func buildRouteParamPostIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, posts.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching PostID from request")
		}
		return u
	}
}

// buildRouteParamReactionIconIDFetcher builds a function that fetches a ReactionIconID from a request routed by chi.
func buildRouteParamReactionIconIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, reactionicons.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ReactionIconID from request")
		}
		return u
	}
}

// buildRouteParamPostReactionIDFetcher builds a function that fetches a PostReactionID from a request routed by chi.
func buildRouteParamPostReactionIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, postreactions.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching PostReactionID from request")
		}
		return u
	}
}

// buildRouteParamNotificationIDFetcher builds a function that fetches a NotificationID from a request routed by chi.
func buildRouteParamNotificationIDFetcher(logger v1.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, notifications.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching NotificationID from request")
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
		ProvideWebhooksServiceWebhookIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideItemsServiceItemIDFetcher,
		ProvideItemsServiceUserIDFetcher,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with forums app", func(t *testing.T) {
		proj := testprojects.BuildForumsApp()
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
		ProvideWebhooksServiceWebhookIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideForumsServiceForumIDFetcher,
		ProvideSubforumsServiceForumIDFetcher,
		ProvideSubforumsServiceSubforumIDFetcher,
		ProvideThreadsServiceForumIDFetcher,
		ProvideThreadsServiceSubforumIDFetcher,
		ProvideThreadsServiceThreadIDFetcher,
		ProvideThreadsServiceUserIDFetcher,
		ProvidePostsServiceForumIDFetcher,
		ProvidePostsServiceSubforumIDFetcher,
		ProvidePostsServiceThreadIDFetcher,
		ProvidePostsServicePostIDFetcher,
		ProvidePostsServiceUserIDFetcher,
		ProvideReactionIconsServiceReactionIconIDFetcher,
		ProvidePostReactionsServiceForumIDFetcher,
		ProvidePostReactionsServiceSubforumIDFetcher,
		ProvidePostReactionsServiceThreadIDFetcher,
		ProvidePostReactionsServicePostIDFetcher,
		ProvidePostReactionsServicePostReactionIDFetcher,
		ProvidePostReactionsServiceUserIDFetcher,
		ProvideNotificationsServiceNotificationIDFetcher,
		ProvideNotificationsServiceUserIDFetcher,
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

		owner := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Thing"),
		}
		typ := models.DataType{
			Name: wordsmith.FromSingularPascalCase("AnotherThing"),
		}

		x := buildProvideSomethingServiceThingIDFetcher(proj, typ, owner)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	anotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/anotherthings"
)

// ProvideAnotherThingsServiceThingIDFetcher provides a ThingIDFetcher.
func ProvideAnotherThingsServiceThingIDFetcher(logger v1.Logger) anotherthings.ThingIDFetcher {
	return buildRouteParamThingIDFetcher(logger)
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
