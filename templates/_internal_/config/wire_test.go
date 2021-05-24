package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

// BEGIN it'd be neat if wire could do this for me one day.

// ProvideConfigServerSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigServerSettings(c *ServerConfig) ServerSettings {
	return c.Server
}

// ProvideConfigAuthSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigAuthSettings(c *ServerConfig) AuthSettings {
	return c.Auth
}

// ProvideConfigDatabaseSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigDatabaseSettings(c *ServerConfig) DatabaseSettings {
	return c.Database
}

// ProvideConfigFrontendSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigFrontendSettings(c *ServerConfig) FrontendSettings {
	return c.Frontend
}

// ProvideSearchSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideSearchSettings(c *ServerConfig) SearchSettings {
	return c.Search
}

// END it'd be neat if wire could do this for me one day.

var (
	// Providers represents this package's offering to the dependency manager.
	Providers = wire.NewSet(
		ProvideConfigServerSettings,
		ProvideConfigAuthSettings,
		ProvideConfigDatabaseSettings,
		ProvideConfigFrontendSettings,
		ProvideSearchSettings,
		ProvideSessionManager,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideConfigServerSettings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideConfigServerSettings()

		expected := `
package example

import ()

// ProvideConfigServerSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigServerSettings(c *ServerConfig) ServerSettings {
	return c.Server
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideConfigAuthSettings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideConfigAuthSettings()

		expected := `
package example

import ()

// ProvideConfigAuthSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigAuthSettings(c *ServerConfig) AuthSettings {
	return c.Auth
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideConfigDatabaseSettings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideConfigDatabaseSettings()

		expected := `
package example

import ()

// ProvideConfigDatabaseSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigDatabaseSettings(c *ServerConfig) DatabaseSettings {
	return c.Database
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideConfigFrontendSettings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideConfigFrontendSettings()

		expected := `
package example

import ()

// ProvideConfigFrontendSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideConfigFrontendSettings(c *ServerConfig) FrontendSettings {
	return c.Frontend
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideSearchSettings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideSearchSettings()

		expected := `
package example

import ()

// ProvideSearchSettings is an obligatory function that
// we're required to have because wire doesn't do it for us.
func ProvideSearchSettings(c *ServerConfig) SearchSettings {
	return c.Search
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProviders(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWireProviders(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers represents this package's offering to the dependency manager.
	Providers = wire.NewSet(
		ProvideConfigServerSettings,
		ProvideConfigAuthSettings,
		ProvideConfigDatabaseSettings,
		ProvideConfigFrontendSettings,
		ProvideSearchSettings,
		ProvideSessionManager,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
