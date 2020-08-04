package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mainTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mainTestDotGo(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	selenium "github.com/tebeka/selenium"
	"testing"
)

func runTestOnAllSupportedBrowsers(t *testing.T, tp testProvider) {
	for _, bn := range []string{"firefox", "chrome"} {
		browserName := bn
		caps := selenium.Capabilities{"browserName": browserName}
		wd, err := selenium.NewRemote(caps, seleniumHubAddr)
		if err != nil {
			panic(err)
		}

		t.Run(bn, tp(wd))
		assert.NoError(t, wd.Quit())
	}
}

type testProvider func(driver selenium.WebDriver) func(t *testing.T)

func TestLoginPage(T *testing.T) {
	runTestOnAllSupportedBrowsers(T, func(driver selenium.WebDriver) func(t *testing.T) {
		return func(t *testing.T) {
			// Navigate to the login page.
			require.NoError(t, driver.Get(urlToUse+"/login"))

			// fetch the button.
			elem, err := driver.FindElement(selenium.ByID, "loginButton")
			if err != nil {
				panic(err)
			}

			// check that it is visible.
			actual, err := elem.IsDisplayed()
			assert.NoError(t, err)
			assert.True(t, actual)
		}
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRunTestOnAllSupportedBrowsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRunTestOnAllSupportedBrowsers()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	selenium "github.com/tebeka/selenium"
	"testing"
)

func runTestOnAllSupportedBrowsers(t *testing.T, tp testProvider) {
	for _, bn := range []string{"firefox", "chrome"} {
		browserName := bn
		caps := selenium.Capabilities{"browserName": browserName}
		wd, err := selenium.NewRemote(caps, seleniumHubAddr)
		if err != nil {
			panic(err)
		}

		t.Run(bn, tp(wd))
		assert.NoError(t, wd.Quit())
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvider(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvider()

		expected := `
package example

import (
	selenium "github.com/tebeka/selenium"
	"testing"
)

type testProvider func(driver selenium.WebDriver) func(t *testing.T)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestLoginPage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestLoginPage()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	selenium "github.com/tebeka/selenium"
	"testing"
)

func TestLoginPage(T *testing.T) {
	runTestOnAllSupportedBrowsers(T, func(driver selenium.WebDriver) func(t *testing.T) {
		return func(t *testing.T) {
			// Navigate to the login page.
			require.NoError(t, driver.Get(urlToUse+"/login"))

			// fetch the button.
			elem, err := driver.FindElement(selenium.ByID, "loginButton")
			if err != nil {
				panic(err)
			}

			// check that it is visible.
			actual, err := elem.IsDisplayed()
			assert.NoError(t, err)
			assert.True(t, actual)
		}
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
