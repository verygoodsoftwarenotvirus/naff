package httpserver

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items"
	"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients"
	"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users"
	"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks"
)

func TestProvideUserIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUserIDFetcher()
	})
}

func TestProvideItemIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideItemIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideUsernameFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUsernameFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideAuthUserIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideAuthUserIDFetcher()
	})
}

func TestProvideWebhooksUserIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksUserIDFetcher()
	})
}

func TestProvideWebhookIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhookIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideOAuth2ServiceClientIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideOAuth2ServiceClientIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestUserIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("obligatory", func(t *testing.T) {
		req := buildRequest(t)
		expected := uint64(123)
		req = req.WithContext(context.WithValue(req.Context(), models.UserIDKey, expected))
		actual := UserIDFetcher(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildChiUserIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("happy path", func(t *testing.T) {
		fn := buildChiUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					users.URIParamKey,
				},
				Values: []string{
					fmt.Sprintf("%d", expected),
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
	T.Run("with invalid value somehow", func(t *testing.T) {
		fn := buildChiUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					users.URIParamKey,
				},
				Values: []string{
					"expected",
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildChiItemIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("happy path", func(t *testing.T) {
		fn := buildChiItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					items.URIParamKey,
				},
				Values: []string{
					fmt.Sprintf("%d", expected),
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
	T.Run("with invalid value somehow", func(t *testing.T) {
		fn := buildChiItemIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					items.URIParamKey,
				},
				Values: []string{
					"expected",
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildChiWebhookIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("happy path", func(t *testing.T) {
		fn := buildChiWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					webhooks.URIParamKey,
				},
				Values: []string{
					fmt.Sprintf("%d", expected),
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
	T.Run("with invalid value somehow", func(t *testing.T) {
		fn := buildChiWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					webhooks.URIParamKey,
				},
				Values: []string{
					"expected",
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildChiOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()
	T.Run("happy path", func(t *testing.T) {
		fn := buildChiOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					oauth2clients.URIParamKey,
				},
				Values: []string{
					fmt.Sprintf("%d", expected),
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
	T.Run("with invalid value somehow", func(t *testing.T) {
		fn := buildChiOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)
		req := buildRequest(t)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					oauth2clients.URIParamKey,
				},
				Values: []string{
					"expected",
				},
			},
		}))
		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}