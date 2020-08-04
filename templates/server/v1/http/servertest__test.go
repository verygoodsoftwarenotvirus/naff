package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_serverTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := serverTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"testing"
)

func buildTestServer() *Server {
	s := &Server{
		DebugMode:  true,
		db:         v1.BuildMockDatabase(),
		config:     &config.ServerConfig{},
		encoder:    &mock.EncoderDecoder{},
		httpServer: provideHTTPServer(),
		logger:     noop.ProvideNoopLogger(),
		frontendService: frontend.ProvideFrontendService(
			noop.ProvideNoopLogger(),
			config.FrontendSettings{},
		),
		webhooksService:      &mock1.WebhookDataServer{},
		usersService:         &mock1.UserDataServer{},
		authService:          &auth.Service{},
		itemsService:         &mock1.ItemDataServer{},
		oauth2ClientsService: &mock1.OAuth2ClientDataServer{},
	}

	return s
}

func TestProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock2.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock2.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid cookie secret", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock2.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISSTRINGISNTLONGENOUGH:(",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock2.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching webhooks", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock2.Anything).Return((*v11.WebhookList)(nil), errors.New("blah"))

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock2.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideServerArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		cookieSecret := "COOKIE_SECRET"
		x := buildProvideServerArgs(proj, cookieSecret)

		expected := `
package main

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

func main() {
	exampleFunction(
		ctx,
		&config.ServerConfig{
			Auth: config.AuthSettings{
				CookieSecret: "COOKIE_SECRET",
			},
		},
		&auth.Service{},
		&frontend.Service{},
		&items.Service{},
		&users.Service{},
		&oauth2clients.Service{},
		&webhooks.Service{},
		mockDB,
		noop.ProvideNoopLogger(),
		&mock.EncoderDecoder{},
		newsman.NewNewsman(nil, nil),
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTestServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildTestServer(proj)

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
)

func buildTestServer() *Server {
	s := &Server{
		DebugMode:  true,
		db:         v1.BuildMockDatabase(),
		config:     &config.ServerConfig{},
		encoder:    &mock.EncoderDecoder{},
		httpServer: provideHTTPServer(),
		logger:     noop.ProvideNoopLogger(),
		frontendService: frontend.ProvideFrontendService(
			noop.ProvideNoopLogger(),
			config.FrontendSettings{},
		),
		webhooksService:      &mock1.WebhookDataServer{},
		usersService:         &mock1.UserDataServer{},
		authService:          &auth.Service{},
		itemsService:         &mock1.ItemDataServer{},
		oauth2ClientsService: &mock1.OAuth2ClientDataServer{},
	}

	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestProvideServer(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"testing"
)

func TestProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock1.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid cookie secret", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISSTRINGISNTLONGENOUGH:(",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock1.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching webhooks", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return((*v11.WebhookList)(nil), errors.New("blah"))

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&items.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mock1.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
