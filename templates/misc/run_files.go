package misc

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func runServerDotXML(dbvendor wordsmith.SuperPalabra) func(*models.Project) string {
	return func(proj *models.Project) string {
		return fmt.Sprintf(`<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="server (%s)" type="GoApplicationRunConfiguration" factoryName="Go Application">
    <module name="%s" />
    <working_directory value="$PROJECT_DIR$" />
    <envs>
      <env name="CONFIGURATION_FILEPATH" value="./environments/testing/config_files/integration-tests-%s.config" />
      <env name="%s_SERVER_LOCAL_CONFIG_STORE_KEY" value="SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU=" />
    </envs>
    <kind value="FILE" />
    <package value="%s" />
    <directory value="$PROJECT_DIR$" />
    <filePath value="$PROJECT_DIR$/cmd/server/main.go" />
    <method v="2" />
  </configuration>
</component>`, dbvendor.RouteName(), proj.Name.RouteName(), dbvendor.RouteName(), strings.ToUpper(proj.Name.Singular()), proj.OutputPath)
	}
}

func runWorkersDotXML(dbvendor wordsmith.SuperPalabra) func(*models.Project) string {
	return func(proj *models.Project) string {
		return fmt.Sprintf(`<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="workers (%s)" type="GoApplicationRunConfiguration" factoryName="Go Application">
    <module name="%s" />
    <working_directory value="$PROJECT_DIR$" />
    <envs>
      <env name="CONFIGURATION_FILEPATH" value="./environments/testing/config_files/integration-tests-%s.config" />
      <env name="%s_WORKERS_LOCAL_CONFIG_STORE_KEY" value="SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU=" />
    </envs>
    <kind value="FILE" />
    <package value="%s" />
    <directory value="$PROJECT_DIR$" />
    <filePath value="$PROJECT_DIR$/cmd/workers/main.go" />
    <method v="2" />
  </configuration>
</component>`, dbvendor.RouteName(), proj.Name.RouteName(), dbvendor.RouteName(), strings.ToUpper(proj.Name.Singular()), proj.OutputPath)
	}
}

func internalUnitTestsRunDotXML(proj *models.Project) string {
	return fmt.Sprintf(`<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="unit tests (internal)" type="GoTestRunConfiguration" factoryName="Go Test">
    <module name="%s" />
    <working_directory value="$PROJECT_DIR$/internal" />
    <root_directory value="$PROJECT_DIR$" />
    <kind value="DIRECTORY" />
    <package value="%s" />
    <directory value="$PROJECT_DIR$/internal" />
    <filePath value="$PROJECT_DIR$" />
    <framework value="gotest" />
    <method v="2" />
  </configuration>
</component>`, proj.Name.RouteName(), proj.OutputPath)

}

func pkgUnitTestsRunDotXML(proj *models.Project) string {
	return fmt.Sprintf(`<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="unit tests (pkg)" type="GoTestRunConfiguration" factoryName="Go Test">
    <module name="%s" />
    <working_directory value="$PROJECT_DIR$/pkg" />
    <root_directory value="$PROJECT_DIR$" />
    <kind value="DIRECTORY" />
    <package value="%s" />
    <directory value="$PROJECT_DIR$/pkg" />
    <filePath value="$PROJECT_DIR$" />
    <framework value="gotest" />
    <method v="2" />
  </configuration>
</component>`, proj.Name.RouteName(), proj.OutputPath)
}
