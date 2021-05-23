package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := usersServiceDotGo(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

const (
	serviceName        = "users_service"
	topicName          = "users"
	counterDescription = "number of users managed by the users service"
	counterName        = metrics.CounterName(serviceName)
)

var (
	_ v1.UserDataServer = (*Service)(nil)
)

type (
	// RequestValidator validates request.
	RequestValidator interface {
		Validate(req *http.Request) (bool, error)
	}

	secretGenerator interface {
		GenerateTwoFactorSecret() (string, error)
		GenerateSalt() ([]byte, error)
	}

	// UserIDFetcher fetches usernames from requests.
	UserIDFetcher func(*http.Request) uint64

	// Service handles our users.
	Service struct {
		cookieSecret        []byte
		userDataManager     v1.UserDataManager
		authenticator       auth.Authenticator
		logger              v11.Logger
		encoderDecoder      encoding.EncoderDecoder
		userIDFetcher       UserIDFetcher
		userCounter         metrics.UnitCounter
		reporter            newsman.Reporter
		secretGenerator     secretGenerator
		userCreationEnabled bool
	}
)

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	authSettings config.AuthSettings,
	logger v11.Logger,
	userDataManager v1.UserDataManager,
	authenticator auth.Authenticator,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	if userIDFetcher == nil {
		return nil, errors.New("userIDFetcher must be provided")
	}

	counter, err := counterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		cookieSecret:        []byte(authSettings.CookieSecret),
		logger:              logger.WithName(serviceName),
		userDataManager:     userDataManager,
		authenticator:       authenticator,
		userIDFetcher:       userIDFetcher,
		encoderDecoder:      encoder,
		userCounter:         counter,
		reporter:            reporter,
		secretGenerator:     &standardSecretGenerator{},
		userCreationEnabled: authSettings.EnableUserSignup,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersServiceConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUsersServiceConstDefs(proj)

		expected := `
package example

import (
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
)

const (
	serviceName        = "users_service"
	topicName          = "users"
	counterDescription = "number of users managed by the users service"
	counterName        = metrics.CounterName(serviceName)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersServiceVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUsersServiceVarDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var (
	_ v1.UserDataServer = (*Service)(nil)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersServiceTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUsersServiceTypeDefs(proj)

		expected := `
package example

import (
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	// RequestValidator validates request.
	RequestValidator interface {
		Validate(req *http.Request) (bool, error)
	}

	secretGenerator interface {
		GenerateTwoFactorSecret() (string, error)
		GenerateSalt() ([]byte, error)
	}

	// UserIDFetcher fetches usernames from requests.
	UserIDFetcher func(*http.Request) uint64

	// Service handles our users.
	Service struct {
		cookieSecret        []byte
		userDataManager     v1.UserDataManager
		authenticator       auth.Authenticator
		logger              v11.Logger
		encoderDecoder      encoding.EncoderDecoder
		userIDFetcher       UserIDFetcher
		userCounter         metrics.UnitCounter
		reporter            newsman.Reporter
		secretGenerator     secretGenerator
		userCreationEnabled bool
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideUsersService(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	authSettings config.AuthSettings,
	logger v1.Logger,
	userDataManager v11.UserDataManager,
	authenticator auth.Authenticator,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	if userIDFetcher == nil {
		return nil, errors.New("userIDFetcher must be provided")
	}

	counter, err := counterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		cookieSecret:        []byte(authSettings.CookieSecret),
		logger:              logger.WithName(serviceName),
		userDataManager:     userDataManager,
		authenticator:       authenticator,
		userIDFetcher:       userIDFetcher,
		encoderDecoder:      encoder,
		userCounter:         counter,
		reporter:            reporter,
		secretGenerator:     &standardSecretGenerator{},
		userCreationEnabled: authSettings.EnableUserSignup,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
