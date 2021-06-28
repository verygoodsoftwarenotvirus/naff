package misc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]func(*models.Project) string{
		".gitignore":           gitIgnore,
		"Makefile":             makefile,
		".gitlab-ci.yml":       gitlabCIDotYAML,
		"README.md":            readmeDotMD,
		".golangci.yml":        golancCILintDotYAML,
		"docker_security.rego": dockerSecurityDotRego,
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(project.OutputPath, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		if _, err = f.WriteString(file(project)); err != nil {
			return err
		}
	}

	return nil
}

func dockerSecurityDotRego(_ *models.Project) string {
	return `package main

suspicious_env_keys = [
    "passwd",
    "password",
    "secret",
    "key",
    "access",
    "api_key",
    "apikey",
    "token",
]

pkg_update_commands = [
    "apk upgrade",
    "apt-get upgrade",
    "dist-upgrade",
]

image_tag_list = [
    "latest",
    "LATEST",
]

# Looking for suspicious environemnt variables
deny[msg] {
    input[i].Cmd == "env"
    val := input[i].Value
    contains(lower(val[_]), suspicious_env_keys[_])
    msg = sprintf("Suspicious ENV key found: %s", [val])
}

# Looking for latest docker image used
warn[msg] {
    input[i].Cmd == "from"
    val := split(input[i].Value[0], ":")
    count(val) == 1
    msg = sprintf("Do not use latest tag with image: %s", [val])
}

# Looking for latest docker image used
warn[msg] {
    input[i].Cmd == "from"
    val := split(input[i].Value[0], ":")
    contains(val[1], image_tag_list[_])
    msg = sprintf("Do not use latest tag with image: %s", [input[i].Value])
}

# Looking for apk upgrade command used in Dockerfile
deny[msg] {
    input[i].Cmd == "run"
    val := concat(" ", input[i].Value)
    contains(val, pkg_update_commands[_])
    msg = sprintf("Do not use upgrade commands: %s", [val])
}

# Looking for ADD command instead using COPY command
deny[msg] {
    input[i].Cmd == "add"
    val := concat(" ", input[i].Value)
    msg = sprintf("Use COPY instead of ADD: %s", [val])
}

# sudo usage
deny[msg] {
    input[i].Cmd == "run"
    val := concat(" ", input[i].Value)
    contains(lower(val), "sudo")
    msg = sprintf("Avoid using 'sudo' command: %s", [val])
}`
}

func gitIgnore(project *models.Project) string {
	output := `# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Test binary, build with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# pretty good IDEs
.idea
.vscode/

# Vim
*.swp

# Sqlite databases
*.db

# OSX
.DS_Store

# Go
vendor
mage_output_file.go

# Python
.env
.mypy_cache
__pycache__
artifacts

# frontend things
internal/services/frontend/_vendor_/*

# Log files
*.log

# go test profiles
*.coverprofile
*.profile

`

	if project.SearchEnabled() {
		output += "# bleve indices\n*.bleve\n"
	}

	return output
}

func gitlabCIDotYAML(project *models.Project) string {
	projRoot := project.OutputPath
	projParts := strings.Split(projRoot, "/")

	ciPath := strings.Join([]string{projParts[0], projParts[1]}, "/")
	ciBuildPath := strings.Join([]string{projParts[1], projParts[2]}, "/")

	f := fmt.Sprintf(`stages:
  - quality
  - frontend-testing
  - integration-testing
  - load-testing
  - publish

before_script:
  - mkdir -p /go/src/%s/
  - cp -rf /builds/%s /go/src/%s/
  - cd /go/src/%s

coverage:
  stage: quality
  image: golang:stretch
  variables:
    GOPATH: "/go"
  script:
    - apt-get update -y && apt-get install -y make git gcc musl-dev
    - go run mage.go ensureMage
    - mage vendor coverage

dependency-injection-check:
  stage: quality
  image: golang:stretch
  variables:
    GOPATH: "/go"
  script:
    - apt-get update -y && apt-get install -y make git gcc musl-dev
    - go run mage.go ensureMage
    - mage rewire

golang-format-check:
  stage: quality
  image: golang:stretch
  variables:
    GOPATH: "/go"
  script:
    - apt-get update -y && apt-get install -y make git gcc musl-dev
    - if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi

golang-lint:
  stage: quality
  image: golangci/golangci-lint:latest # v1.18
  variables:
    GO111MODULE: "on"
  script:
    - go mod vendor
    - golangci-lint run --config=.golangci.yml --deadline=15m

#integration_tests_sqlite:
#  stage: integration-testing
#  image: docker/compose:latest
#  services:
#    - docker:dind
#  variables:
#    GOPATH: "/go"
#  script:
#    - ls -Al environments/testing/config_files/integration-tests-sqlite.toml
#    - ls -Al /go/src/%s/environments/testing/config_files/integration-tests-sqlite.toml
#    - whoami
#    - cat /go/src/%s/environments/testing/config_files/integration-tests-sqlite.toml
#    - apk add --no-cache git make musl-dev go
#    - go version
#    - docker-compose --file environments/testing/compose_files/integration_tests/integration-tests-base.yaml --file environments/testing/compose_files/integration_tests/integration-tests-sqlite.yaml up --build --force-recreate --remove-orphans --renew-anon-volumes --always-recreate-deps --abort-on-container-exit

#integration_tests_postgres:
#  stage: integration-testing
#  image: docker/compose:latest
#  services:
#    - docker:dind
#  variables:
#    GOPATH: "/go"
#  script:
#    - apk add --update make
#    - docker-compose --file environments/testing/compose_files/integration_tests/integration-tests-base.yaml --file environments/testing/compose_files/integration_tests/integration-tests-postgres.yaml up --build --force-recreate --remove-orphans --renew-anon-volumes --always-recreate-deps --abort-on-container-exit

#integration_tests_mariadb:
#  stage: integration-testing
#  image: docker/compose:latest
#  services:
#    - docker:dind
#  variables:
#    GOPATH: "/go"
#  script:
#    - apk add --update make
#    - docker-compose --file environments/testing/compose_files/integration_tests/integration-tests-base.yaml --file environments/testing/compose_files/integration_tests/integration-tests-mariadb.yaml up --build --force-recreate --remove-orphans --renew-anon-volumes --always-recreate-deps --abort-on-container-exit

#frontend-selenium-tests:
#  stage: integration-testing
#  image: docker/compose:latest
#  services:
#    - docker:dind
#  script:
#    - apk add --update make
#    - docker-compose --file environments/testing/compose_files/frontend-tests.yaml up --build --force-recreate --remove-orphans --renew-anon-volumes --always-recreate-deps --abort-on-container-exit

# daily load tests

daily-load_tests_postgres:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load_tests_postgres
  only:
    - schedules

daily-load_tests_mariadb:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load_tests_mariadb
  only:
    - schedules

daily-load_tests_sqlite:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load_tests_sqlite
  only:
    - schedules
`, ciPath, ciBuildPath, ciPath, projRoot, projRoot, projRoot)

	return f
}

func golancCILintDotYAML(project *models.Project) string {
	projRoot := project.OutputPath
	f := fmt.Sprintf(`# https://github.com/golangci/golangci-lint/blob/507703b444d95d8c89961bebeedfb22f61cde67c/pkg/config/config.go

# options for analysis running
run:
  # timeout for analysis, e.g. 30s, 5m, default is 1m
  deadline: 5m

  # exit code when at least one issue was found, default is 1
  issues-exit-code: 1

  # include test files or not, default is true
  tests: true

  # list of build tags, all linters use it. Default is empty list.
  build-tags:
    - mytag

  # which dirs to skip: they won't be analyzed;
  # can use regexp here: generated.*, regexp is applied on full path;
  # default value is empty list, but next dirs are always skipped independently
  # from this option's value:
  #   	vendor$, third_party$, testdata$, examples$, Godeps$, builtin$
  skip-dirs:
    - cmd/tools
    - cmd/playground

  # which files to skip: they will be analyzed, but issues from them
  # won't be reported. Default value is empty list, but there is
  # no need to include all autogenerated files, we confidently recognize
  # autogenerated files. If it's not please let us know.
  # skip-files:
  #   -

  # by default isn't set. If set we pass it to "go list -mod={option}". From "go help modules":
  # If invoked with -mod=readonly, the go command is disallowed from the implicit
  # automatic updating of go.mod described above. Instead, it fails when any changes
  # to go.mod are needed. This setting is most useful to check that go.mod does
  # not need updates, such as in a continuous integration and testing system.
  # If invoked with -mod=vendor, the go command assumes that the vendor
  # directory holds the correct copies of dependencies and ignores
  # the dependency descriptions in go.mod.
  #
  # available options: readonly|vendor
  # modules-download-mode: release

# output configuration options
output:
  # colored-line-number|line-number|json|tab|checkstyle|code-climate
  format: colored-line-number

  # print lines of code with issue, default is true
  print-issued-lines: true

  # print linter name in the end of issue text, default is true
  print-linter-name: true

# all available settings of specific linters
linters-settings:
  cyclop:
    # the maximal code complexity to report
    max-complexity: 25
    # the maximal average package complexity. If it's higher than 0.0 (float) the check is enabled (default 0.0)
    package-average: 0.0
    # should ignore tests (default false)
    skip-tests: true
  errcheck:
    # report about not checking of errs in type assertions: `+"`"+`a := b.(MyStruct)`+"`"+`;
    # default is false: such cases aren't reported by default.
    check-type-assertions: true

    # report about assignment of errs to blank identifier: `+"`"+`num, _ := strconv.Atoi(numStr)`+"`"+`;
    # default is false: such cases aren't reported by default.
    check-blank: true

    # # path to a file containing a list of functions to exclude from checking
    # # see https://github.com/kisielk/errcheck#excluding-functions for details
    # exclude: /path/to/file.txt
  forbidigo:
    # Forbid the following identifiers
    forbid:
      - ^t\.SkipNow\(\)$ # no skipped tests
    # Exclude godoc examples from forbidigo checks.  Default is true.
    exclude_godoc_examples: false
  govet:
    # report about shadowed variables
    check-shadowing: true

    # settings per analyzer
    settings:
      printf: # analyzer name, run `+"`"+`go tool vet help`+"`"+` to see all analyzers
        funcs: # run `+"`"+`go tool vet help printf`+"`"+` to see available settings for `+"`"+`printf`+"`"+` analyzer
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf
    # enable or disable analyzers by name
    enable-all: true
#  goheader:
#    values:
#      const:
#      # define here const type values in format k:v, for example:
#      # COMPANY: MY COMPANY
#      regexp:
#      # define here regexp type values, for example
#      # AUTHOR: .*@mycompany\.com
#    template-path: development/header_template.md
  gofmt:
    # simplify code: gofmt with `+"`"+`-s`+"`"+` option, true by default
    simplify: true
  goimports:
    # put imports beginning with prefix after 3rd-party packages;
    # it's a comma-separated list of prefixes
    local-prefixes: %s
  gocyclo:
    # minimal code complexity to report, 30 by default (but we recommend 10-20)
    min-complexity: 20
  gosec:
    exclude:
      - G304
  maligned:
    # print struct with more effective memory layout or not, false by default
    suggest-new: true
  dupl:
    # tokens count to trigger issue, 150 by default
    threshold: 100
  goconst:
    # minimal length of string constant, 3 by default
    min-len: 3
    # minimal occurrences count to trigger, 3 by default
    min-occurrences: 3
  lll:
    # max line length, lines longer will be reported. Default is 120.
    # '\t' is counted as 1 character by default, and can be changed with the tab-width option
    line-length: 512
    # tab width in spaces. Default to 1.
    tab-width: 1
  misspell:
    # Correct spellings using locale preferences for US or UK.
    # Default is to use a neutral variety of English.
    # Setting locale to US will correct the British spelling of 'colour' to 'color'.
    locale: US
    # ignore-words:
    #   - someword
  nestif:
    min-complexity: 8
  wsl:
    allow-cuddle-declarations: true
  unused:
    # treat code as a program (not a library) and report unused exported identifiers; default is false.
    # XXX: if you enable this setting, unused will report a lot of false-positives in text editors:
    # if it's called for subdir of a project it can't find funcs usages. All text editor integrations
    # with golangci-lint call it on a directory with the changed file.
    check-exported: false
  unparam:
    # Inspect exported functions, default is false. Set to true if no external program/library imports your code.
    # XXX: if you enable this setting, unparam will report a lot of false-positives in text editors:
    # if it's called for subdir of a project it can't find external interfaces. All text editor integrations
    # with golangci-lint call it on a directory with the changed file.
    check-exported: false
  nakedret:
    # make an issue if func has more lines of code than this setting and it has naked returns; default is 30
    max-func-lines: 4
  prealloc:
    # XXX: we don't recommend using this linter before doing performance profiling.
    # For most programs usage of prealloc will be a premature optimization.

    # Report preallocation suggestions only on simple loops that have no returns/breaks/continues/gotos in them.
    # True by default.
    simple: true
    range-loops: true # Report preallocation suggestions on range loops, true by default
    for-loops: false # Report preallocation suggestions on for loops, false by default
  gocritic:
    # Which checks should be disabled; can't be combined with 'enabled-checks'; default is empty
    disabled-checks:
      - captLocal
      - singleCaseSwitch

    # Enable multiple checks by tags, run `+"`"+`GL_DEBUG=gocritic golangci-lint`+"`"+` run to see all tags and checks.
    # Empty list by default. See https://github.com/go-critic/go-critic#usage -> section "Tags".
    enabled-tags:
      - diagnostic
      - style
      - performance
      - opinionated

    settings: # settings passed to gocritic
      rangeValCopy:
        sizeThreshold: 32

# last updated this list at v1.41.0
linters:
  fast: false
  disable-all: false
  enable:
    # Enabled By Default
    - deadcode # Finds unused code
    - errcheck # Errcheck is a program for checking for unchecked errs in go programs. These unchecked errs can be critical bugs in some cases
    - gosimple # Linter for Go source code that specializes in simplifying a code
    - govet # Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
    - ineffassign # Detects when assignments to existing variables are not used
    - staticcheck # Staticcheck is a go vet on steroids, applying a ton of static analysis checks
    - structcheck # Finds unused struct fields
    - typecheck # Like the front-end of a Go compiler, parses and type-checks Go code
    - unused # Checks Go code for unused constants, variables, functions and types
    - varcheck # Finds unused global variables and constants
    # Disabled By Default
    - asciicheck       # Simple linter to check that your code does not contain non-ASCII identifiers
    - bodyclose        # checks whether HTTP response body is closed successfully
    - cyclop           # checks function and package cyclomatic complexity
    - depguard         # Go linter that checks if package imports are in a list of acceptable packages
    - dogsled          # Checks assignments with too many blank identifiers (e.g. x, , , _, := f())
    # - dupl             # Tool for code clone detection
    - durationcheck    # check for two durations multiplied together
    - errorlint        # errorlint is a linter for that can be used to find code that will cause problems with the error wrapping scheme introduced in Go 1.13.
    - exhaustive       # check exhaustiveness of enum switch statementa
    # - exhaustivestruct # Checks if all struct's fields are initialized
    - exportloopref    # checks for pointers to enclosing loop variables
    - forbidigo        # Forbids identifiers
    - forcetypeassert  # finds forced type assertions
    # - funlen           # Tool for detection of long functions
    # - gci              # Gci control golang package import order and make it always deterministic.
    # - gochecknoglobals # check that no global variables exist
    # - gochecknoinits   # Checks that no init functions are present in Go code
    - gocognit         # Computes and checks the cognitive complexity of functions
    - goconst          # Finds repeated strings that could be replaced by a constant
    - gocritic         # Provides many diagnostics that check for bugs, performance and style issues.
    - gocyclo          # Computes and checks the cyclomatic complexity of functions
    - godot            # Check if comments end in a period
    - godox            # Tool for detection of FIXME, TODO and other comment keywords
    - goerr113         # Golang linter to check the errors handling expressions
    - gofmt            # Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
    # - gofumpt          # Gofumpt checks whether code was gofumpt-ed.
    - goheader         # Checks is file header matches to pattern
    - goimports        # Goimports does everything that gofmt does. Additionally it checks unused imports
    - gomnd            # An analyzer to detect magic numbers.
    - gomoddirectives  # Manage the use of 'replace', 'retract', and 'excludes' directives in go.mod.
    - gomodguard       # Allow and block list linter for direct Go module dependencies. This is different from depguard where there are different block types for example version constraints and module recommendations.
    - goprintffuncname # Checks that printf-like functions are named with f at the end
    - gosec            # Inspects source code for security problems
    - ifshort          # Checks that your code uses short syntax for if-statements whenever possible
    - importas         # Enforces consistent import aliases
    - makezero         # Finds slice declarations with non-zero initial length
    - misspell         # Finds commonly misspelled English words in comments
    - nakedret         # Finds naked returns in functions greater than a specified function length
    - nestif           # Reports deeply nested if statements	complexity
    - nilerr           # Finds the code that returns nil even if it checks that the error is not nil.
    # - nlreturn         # nlreturn checks for a new line before return and branch statements to increase code clarity
    - noctx            # noctx finds sending http request without context.Context
    - nolintlint       # Reports ill-formed or insufficient nolint directives
    - paralleltest     # paralleltest detects missing usage of t.Parallel() method in your Go test
    - prealloc         # Finds slice declarations that could potentially be preallocated
    - predeclared      # find code that shadows one of Go's predeclared identifiers
    - promlinter       # Check Prometheus metrics naming via promlint
    - revive           # Fast, configurable, extensible, flexible, and beautiful linter for Go. Drop-in replacement of golint.
    # - rowserrcheck     # checks whether Err of rows is checked successfully
    - sqlclosecheck    # Checks that sql.Rows and sql.Stmt are closed.
    # - stylecheck       # Stylecheck is a replacement for golint	style
    # - tagliatelle      # Checks the struct tags.
    # - testpackage      # linter that makes you use a separate _test package
    - thelper          # thelper detects golang test helpers without t.Helper() call and checks the consistency of test helpers
    - tparallel        # tparallel detects inappropriate usage of t.Parallel() method in your Go test codes
    - unconvert        # Remove unnecessary type conversions
    - unparam          # Reports unused function parameters	unused
    - wastedassign     # wastedassign finds wasted assignment statements.
    - whitespace       # Tool for detection of leading and trailing whitespace
    # - wrapcheck        # Checks that errors returned from external packages are wrapped
    # - wsl              # Whitespace Linter - Forces you to use empty lines!

  disable:
    - exhaustivestruct # Checks if all struct's fields are initialized
    - gci              # control package import order and make it always deterministic.
    - gochecknoinits   # Checks that no init functions are present in Go code
    - gochecknoglobals # Checks that no globals are present in Go code
    - dupl             # Tool for code clone detection
    - funlen           # Tool for detection of long functions
    - gofumpt          # Gofumpt checks whether code was gofumpt-ed.
    - lll              # Reports long lines
    - nlreturn         # nlreturn checks for a new line before return and branch statements to increase code clarity
    - rowserrcheck     # checks whether Err of rows is checked successfully; lots of false positives
    - tagliatelle      # Checks the struct tags.
    - stylecheck       # Stylecheck is a replacement for golint
    - testpackage      # linter that makes you use a separate _test package
    - wsl              # Whitespace Linter - Forces you to use empty lines! Easily the most annoying one on here
    - wrapcheck        # Checks that errs returned from external packages are wrapped

issues:
  # List of regexps of issue texts to exclude, empty list by default.
  # But independently from this option we use default exclude patterns,
  # it can be disabled by `+"`"+`exclude-use-default: false`+"`"+`. To list all
  # excluded by default patterns execute `+"`"+`golangci-lint run --help`+"`"+`
  # exclude:
  #   -

  # Excluding configuration per-path, per-linter, per-text and per-source
  exclude-rules:
    # Exclude some linters from running on tests files.
    - path: _test\.go
      linters:
        - gocyclo
        - goconst # I want my tests to repeat themselves
        - errcheck
        - dupl
        - gosec
        - lll
        - goerr113
        - bodyclose
        - unparam
        - gocognit
        - gomnd

    - path: tests/
      linters:
        - gosec
        - gomnd

    - path: pkg/client/httpclient/roundtripper_*
      linters:
        - bodyclose

    - path: tests/testutil/
      linters:
        - bodyclose
        - gocyclo

    - path: mock_.*\.go
      linters:
        - lll
        - mnd

    - path: pkg/types/
      linters:
        - gocyclo # the update funcs can have very high cyclomatic complexities

    - path: pkg/types/fakes
      linters:
        - gomnd

    - path: pkg/types/mock
      linters:
        - gomnd

    - path: internal/audit/
      linters:
        - gosec

    - path: internal/services/
      linters:
        - wsl
      text: "return statements should not be cuddled if block has more than two lines"

    ## BEGIN SPECIAL SNOWFLAKES

    - path: internal/services/frontend/time_test.go
      linters:
        - paralleltest

    - path: tests/frontend
      linters:
        - deadcode
        - unused
        - paralleltest
        - thelper

    - path: cmd/config_gen
      linters:
        - gomnd

    - path: internal/database/querybuilding/mariadb/migrations.go
      linters:
        - gomnd

    - path: internal/database/querybuilding/postgres/migrations.go
      linters:
        - gomnd

    - path: internal/database/querybuilding/sqlite/migrations.go
      linters:
        - gomnd

    - path: internal/permissions/permissions.go
      linters:
        - deadcode
        - varcheck

    - path: tests/load/actions.go
      linters:
        - gocyclo

    - path: internal/observability/logging/zap
      linters:
        - gocyclo

    - path: tests/load/main.go
      linters:
        - gocritic # gocritic complains about an interface implementation I have no control over

    - path: internal/server/httpclient/server.go
      linters:
        - maligned # these structs are never copied and are structured for documentation purposes

    - path: cmd/config_gen/main.go
      linters:
        - gosec
        - lll

    - path: internal/config/config.go
      linters:
        - lll

    - path: tests/testutil/testutil.go
      linters:
        - goerr113

    - path: cmd/server/main.go
      linters:
        - goerr113

    - path: internal/database/querybuilding/postgres/migrations.go
      linters:
        - gocognit
        - gocyclo

    ## END SPECIAL SNOWFLAKES

    # Exclude known linters from partially hard-vendored code, which is impossible to exclude via "nolint" comments.
    # - path: internal/hmac/
    #   text: "weak cryptographic primitive"
    #   linters:
    #     - gosec

    - linters:
        - gocritic
      text: "appendAssign: "

    - linters:
        - goerr113
      text: "do not define dynamic errs"

    # ignore this error type because it isn't defined anywhere, and it's detecting a false positive
    - linters:
        - gosec
      text: "G304:"

    # Exclude lll issues for long lines with go:generate
    - linters:
        - lll
      source: "^//go:generate "

  # Independently from option `+"`"+`exclude`+"`"+` we use default exclude patterns, it can be disabled by this option. To list all
  # excluded by default patterns execute `+"`"+`golangci-lint run --help`+"`"+`.
  # Default value for this option is true.
  exclude-use-default: false

  # Maximum issues count per one linter. Set to 0 to disable. Default is 50.
  max-issues-per-linter: 0

  # Maximum count of issues with the same text. Set to 0 to disable. Default is 3.
  max-same-issues: 0

  # Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only
  # changes in HEAD~ are analyzed. It's a super-useful option for integration of golangci-lint into existing large
  # codebase. It's not practical to fix all existing issues at the moment of integration: much better don't allow issues
  # in new code. Default is false.
  new: false

  #
  # # Show only new issues created after git revision `+"`"+`REV`+"`"+`
  # new-from-rev: REV
  #
  # # Show only new issues created in git patch with set file path.
  # new-from-patch: path/to/patch/file
  #`, projRoot)

	return f
}
