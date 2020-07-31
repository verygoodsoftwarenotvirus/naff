package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_configTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := configTestDotGo(proj)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"io/ioutil"
	"os"
	"testing"
)

func Test_randString(t *testing.T) {
	t.Parallel()

	actual := randString(randStringSize)
	assert.NotEmpty(t, actual)
	assert.Len(t, actual, 52)
}

func TestBuildConfig(t *testing.T) {
	t.Parallel()

	actual := BuildConfig()
	assert.NotNil(t, actual)
}

func TestParseConfigFile(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		tf, err := ioutil.TempFile(os.TempDir(), "*.toml")
		require.NoError(t, err)
		expected := "thisisatest"

		_, err = tf.Write([]byte(fmt.Sprintf(`+"`"+`
[server]
http_port = 1234
debug = false

[database]
provider = "postgres"
debug = true
connection_details = "%s"
`+"`"+`, expected)))
		require.NoError(t, err)

		expectedConfig := &ServerConfig{
			Server: ServerSettings{
				HTTPPort: 1234,
				Debug:    false,
			},
			Database: DatabaseSettings{
				Provider:          "postgres",
				Debug:             true,
				ConnectionDetails: v1.ConnectionDetails(expected),
			},
		}

		cfg, err := ParseConfigFile(tf.Name())
		assert.NoError(t, err)

		assert.Equal(t, expectedConfig.Server.HTTPPort, cfg.Server.HTTPPort)
		assert.Equal(t, expectedConfig.Server.Debug, cfg.Server.Debug)
		assert.Equal(t, expectedConfig.Database.Provider, cfg.Database.Provider)
		assert.Equal(t, expectedConfig.Database.Debug, cfg.Database.Debug)
		assert.Equal(t, expectedConfig.Database.ConnectionDetails, cfg.Database.ConnectionDetails)

		assert.NoError(t, os.Remove(tf.Name()))
	})

	T.Run("with nonexistent file", func(t *testing.T) {
		cfg, err := ParseConfigFile("/this/doesn't/even/exist/lol")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_randString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTest_randString()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_randString(t *testing.T) {
	t.Parallel()

	actual := randString(randStringSize)
	assert.NotEmpty(t, actual)
	assert.Len(t, actual, 52)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBuildConfig(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestBuildConfig()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	t.Parallel()

	actual := BuildConfig()
	assert.NotNil(t, actual)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestParseConfigFile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestParseConfigFile(proj)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigFile(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		tf, err := ioutil.TempFile(os.TempDir(), "*.toml")
		require.NoError(t, err)
		expected := "thisisatest"

		_, err = tf.Write([]byte(fmt.Sprintf(`+"`"+`
[server]
http_port = 1234
debug = false

[database]
provider = "postgres"
debug = true
connection_details = "%s"
`+"`"+`, expected)))
		require.NoError(t, err)

		expectedConfig := &ServerConfig{
			Server: ServerSettings{
				HTTPPort: 1234,
				Debug:    false,
			},
			Database: DatabaseSettings{
				Provider:          "postgres",
				Debug:             true,
				ConnectionDetails: v1.ConnectionDetails(expected),
			},
		}

		cfg, err := ParseConfigFile(tf.Name())
		assert.NoError(t, err)

		assert.Equal(t, expectedConfig.Server.HTTPPort, cfg.Server.HTTPPort)
		assert.Equal(t, expectedConfig.Server.Debug, cfg.Server.Debug)
		assert.Equal(t, expectedConfig.Database.Provider, cfg.Database.Provider)
		assert.Equal(t, expectedConfig.Database.Debug, cfg.Database.Debug)
		assert.Equal(t, expectedConfig.Database.ConnectionDetails, cfg.Database.ConnectionDetails)

		assert.NoError(t, os.Remove(tf.Name()))
	})

	T.Run("with nonexistent file", func(t *testing.T) {
		cfg, err := ParseConfigFile("/this/doesn't/even/exist/lol")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
