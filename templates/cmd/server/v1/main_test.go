package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"os"
)

func main() {
	// initialize our logger of choice.
	logger := zerolog.NewZeroLogger()

	// find and validate our configuration filepath.
	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")
	if configFilepath == "" {
		logger.Fatal(errors.New("no configuration file provided"))
	}

	// parse our config file.
	cfg, err := config.ParseConfigFile(configFilepath)
	if err != nil || cfg == nil {
		logger.Fatal(fmt.Errorf("error parsing configuration file: %w", err))
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Meta.StartupDeadline)
	ctx, span := tracing.StartSpan(ctx, "initialization")

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := cfg.ProvideDatabaseConnection(logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := cfg.ProvideDatabaseClient(ctx, logger, rawDB)
	if err != nil {
		logger.Fatal(fmt.Errorf("error initializing database client: %w", err))
	}

	// build our server struct.
	logger.Debug("building server")
	server, err := BuildServer(ctx, cfg, logger, dbClient, rawDB)
	span.End()
	cancel()

	if err != nil {
		logger.Fatal(fmt.Errorf("error initializing HTTP server: %w", err))
	}

	// I slept and dreamt that life was joy.
	//   I awoke and saw that life was service.
	//   	I acted and behold, service deployed.
	server.Serve()
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildMain(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"os"
)

func main() {
	// initialize our logger of choice.
	logger := zerolog.NewZeroLogger()

	// find and validate our configuration filepath.
	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")
	if configFilepath == "" {
		logger.Fatal(errors.New("no configuration file provided"))
	}

	// parse our config file.
	cfg, err := config.ParseConfigFile(configFilepath)
	if err != nil || cfg == nil {
		logger.Fatal(fmt.Errorf("error parsing configuration file: %w", err))
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Meta.StartupDeadline)
	ctx, span := tracing.StartSpan(ctx, "initialization")

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := cfg.ProvideDatabaseConnection(logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := cfg.ProvideDatabaseClient(ctx, logger, rawDB)
	if err != nil {
		logger.Fatal(fmt.Errorf("error initializing database client: %w", err))
	}

	// build our server struct.
	logger.Debug("building server")
	server, err := BuildServer(ctx, cfg, logger, dbClient, rawDB)
	span.End()
	cancel()

	if err != nil {
		logger.Fatal(fmt.Errorf("error initializing HTTP server: %w", err))
	}

	// I slept and dreamt that life was joy.
	//   I awoke and saw that life was service.
	//   	I acted and behold, service deployed.
	server.Serve()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
