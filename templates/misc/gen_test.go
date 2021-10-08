package misc

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()
		project.OutputPath = filepath.Join(os.TempDir(), "one", "two")

		assert.NoError(t, RenderPackage(project))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/dev/null`

		assert.Error(t, RenderPackage(proj))
	})
}

func Test_gitIgnore(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()

		expected := `# Binaries for programs and plugins
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

# Python
.env
.mypy_cache
__pycache__
artifacts

# frontend things
node_modules

# Log files
npm-debug.log*
yarn-debug.log*
yarn-error.log*

frontend/v1/public/bundle.*

*.coverprofile
*.profile

cmd/playground
*.bleve
*.sqlite
`
		actual := gitIgnore(project)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_gitlabCIDotYAML(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()

		expected := `stages:
  - quality
  - frontend-testing
  - integration-testing
  - load-testing
  - publish

before_script:
  - mkdir -p /go/src/gitlab.com/verygoodsoftwarenotvirus/
  - cp -rf /builds/verygoodsoftwarenotvirus/naff /go/src/gitlab.com/verygoodsoftwarenotvirus/
  - cd /go/src/gitlab.com/verygoodsoftwarenotvirus/naff/example_output
  - apt-get update -y && apt-get install -y make git gcc musl-dev

formatting:
  stage: quality
  image: golang:stretch
  variables:
    GOPATH: "/go"
  script:
    - apt-get update -y && apt-get install -y make git gcc musl-dev
    - if [ $(gofmt -l . | grep -Ev '^vendor\/' | head -c1 | wc -c) -ne 0 ]; then exit 1; fi

coverage:
  stage: quality
  image: golang:stretch
  variables:
    GOPATH: "/go"
  script:
    - apt-get update -y && apt-get install -y make git gcc musl-dev
    - make coverage

linting:
  stage: quality
  image: golangci/golangci-lint:latest # v1.18
  variables:
    GO111MODULE: "on"
  script:
    - go mod vendor
    - golangci-lint run --config=.golangci.yml --deadline=15m

build-frontend:
  stage: quality
  image: node:10
  before_script:
    - cd frontend/v1
    - npm install
  script:
    - npm run build

integration-tests-postgres:
  stage: integration-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make integration-tests-postgres

integration-tests-mariadb:
  stage: integration-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make integration-tests-mariadb

integration-tests-sqlite:
  stage: integration-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make integration-tests-sqlite

frontend-selenium-tests:
  stage: integration-testing
  image: docker:latest
  services:
    - docker:dind
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make frontend-tests

# load tests
load-tests-postgres:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "2m30s"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-postgres
  except:
    - schedules

load-tests-mariadb:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "5m00s"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-mariadb
  except:
    - schedules

load-tests-sqlite:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "2m30s"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-sqlite
  except:
    - schedules

# daily load tests

daily-load-tests-postgres:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-postgres
  only:
    - schedules

daily-load-tests-mariadb:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-mariadb
  only:
    - schedules

daily-load-tests-sqlite:on-schedule:
  stage: load-testing
  image: docker:latest
  services:
    - docker:dind
  variables:
    GOPATH: "/go"
    LOADTEST_RUN_TIME: "10m"
  script:
    - apk add --update --no-cache py-pip openssl python3-dev libffi-dev
      openssl-dev gcc libc-dev make
    - pip install docker-compose
    - make load-tests-sqlite
  only:
    - schedules

# miscellaneous

#gitlabcr:
#  stage: publish
#  image: docker:latest
#  services:
#    - docker:dind
#  script:
#    - docker login --username=gitlab-ci-token --password=$CI_JOB_TOKEN registry.gitlab.com
#    - docker build --tag registry.gitlab.com/verygoodsoftwarenotvirus/naff/example_output:latest --file environments/production/Dockerfile .
#    - docker push registry.gitlab.com/verygoodsoftwarenotvirus/naff/example_output:latest
#  only:
#    - master
`
		actual := gitlabCIDotYAML(project)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_golancCILintDotYAML(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		project := testprojects.BuildTodoApp()

		expected := `# options for analysis running
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
    # javascript
    - client/v1/frontend
    # borrowed code/utilities
    - cmd/tools

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
  # available options: readonly|release|vendor
  modules-download-mode: vendor

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
  errcheck:
    # report about not checking of errors in type assertions: ` + "`" + `a := b.(MyStruct)` + "`" + `;
    # default is false: such cases aren't reported by default.
    check-type-assertions: true

    # report about assignment of errors to blank identifier: ` + "`" + `num, _ := strconv.Atoi(numStr)` + "`" + `;
    # default is false: such cases aren't reported by default.
    check-blank: true

    # # path to a file containing a list of functions to exclude from checking
    # # see https://github.com/kisielk/errcheck#excluding-functions for details
    # exclude: /path/to/file.txt
  govet:
    # report about shadowed variables
    check-shadowing: true

    # settings per analyzer
    settings:
      printf: # analyzer name, run ` + "`" + `go tool vet help` + "`" + ` to see all analyzers
        funcs: # run ` + "`" + `go tool vet help printf` + "`" + ` to see available settings for ` + "`" + `printf` + "`" + ` analyzer
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Infof
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Warnf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Errorf
          - (github.com/golangci/golangci-lint/pkg/logutils.Log).Fatalf
    # enable or disable analyzers by name
    enable-all: true

  golint:
    # minimal confidence for issues, default is 0.8
    min-confidence: 0.8
  gofmt:
    # simplify code: gofmt with ` + "`" + `-s` + "`" + ` option, true by default
    simplify: true
  goimports:
    # put imports beginning with prefix after 3rd-party packages;
    # it's a comma-separated list of prefixes
    local-prefixes: gitlab.com/verygoodsoftwarenotvirus/naff/example_output
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
    line-length: 251
    # tab width in spaces. Default to 1.
    tab-width: 1
  misspell:
    # Correct spellings using locale preferences for US or UK.
    # Default is to use a neutral variety of English.
    # Setting locale to US will correct the British spelling of 'colour' to 'color'.
    locale: US
    # ignore-words:
    #   - someword
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

    # Enable multiple checks by tags, run ` + "`" + `GL_DEBUG=gocritic golangci-lint` + "`" + ` run to see all tags and checks.
    # Empty list by default. See https://github.com/go-critic/go-critic#usage -> section "Tags".
    enabled-tags:
      - diagnostic
      - style
      - performance
      - opinionated

    settings: # settings passed to gocritic
      rangeValCopy:
        sizeThreshold: 32

linters:
  fast: false
  disable-all: false
  enable:
    # Enabled By Default
    - govet # Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string
    - errcheck # Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases
    - staticcheck # Staticcheck is a go vet on steroids, applying a ton of static analysis checks
    - unused # Checks Go code for unused constants, variables, functions and types
    - gosimple # Linter for Go source code that specializes in simplifying a code
    - structcheck # Finds unused struct fields
    - varcheck # Finds unused global variables and constants
    - ineffassign # Detects when assignments to existing variables are not used
    - deadcode # Finds unused code
    - typecheck # Like the front-end of a Go compiler, parses and type-checks Go code
    # Disabled By Default
    - bodyclose # checks whether HTTP response body is closed successfully
    - golint # Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
    - gosec # Inspects source code for security problems
    - interfacer # Linter that suggests narrower interface types
    - unconvert # Remove unnecessary type conversions
    - goconst # Finds repeated strings that could be replaced by a constant
    - gocyclo # Computes and checks the cyclomatic complexity of functions
    - gofmt # Gofmt checks whether code was gofmt-ed. By default this tool runs with -s option to check for code simplification
    - maligned # Tool to detect Go structs that would take less memory if their fields were sorted
    - depguard # Go linter that checks if package imports are in a list of acceptable packages
    - misspell # Finds commonly misspelled English words in comments
    - lll # Reports long lines
    - unparam # Reports unused function parameters
    - dogsled # Checks assignments with too many blank identifiers (e.g. x, _, _, _, := f())
    - nakedret # Finds naked returns in functions greater than a specified function length
    - prealloc # Finds slice declarations that could potentially be preallocated
    - scopelint # Scopelint checks for unpinned variables in go programs
    - gocritic # The most opinionated Go source code linter
    - godox # Tool for detection of FIXME, TODO and other comment keywords
    - whitespace # Tool for detection of leading and trailing whitespace
    - goimports # Goimports does everything that gofmt does. Additionally it checks unused imports
  disable:
    - stylecheck # Stylecheck is a replacement for golint
    - gochecknoinits # Checks that no init functions are present in Go code
    - gochecknoglobals # Checks that no globals are present in Go code
    - dupl # Tool for code clone detection
    - funlen # Tool for detection of long functions

issues:
  # # List of regexps of issue texts to exclude, empty list by default.
  # # But independently from this option we use default exclude patterns,
  # # it can be disabled by ` + "`" + `exclude-use-default: false` + "`" + `. To list all
  # # excluded by default patterns execute ` + "`" + `golangci-lint run --help` + "`" + `
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
        - bodyclose

    - path: tests/v1
      linters:
        - gosec

    # Exclude routers from gocyclo
    - path: server/v1/http/routes.go
      linters:
        - gocyclo

    # Exclude frontend routers from gocyclo
    - path: services/v1/frontend/http_routes.go
      linters:
        - gocyclo

    - path: tests/v1/testutil/
      linters:
        - bodyclose
        - gocyclo

    - path: mock_.*\.go
      linters:
        - lll

    - path: models/
      linters:
        - gocyclo # the update funcs can have very high cyclomatic complexities

    ## BEGIN SPECIAL SNOWFLAKES

    - path: tests/v1/load/actions.go
      linters:
        - gocyclo

    - path: tests/v1/load/main.go
      linters:
        - gocritic # gocritic complains about an interface implementation I have no control over

    - path: server/v1/http/server.go
      linters:
        - maligned # these structs are never copied and are structured for documentation purposes

    - path: cmd/config_gen/v1/main.go
      linters:
        - gosec
        - lll

    - path: internal/v1/config/config.go
      linters:
        - lll

    ## END SPECIAL SNOWFLAKES

    # Exclude known linters from partially hard-vendored code,
    # which is impossible to exclude via "nolint" comments.
    # - path: internal/hmac/
    #   text: "weak cryptographic primitive"
    #   linters:
    #     - gosec

    # ignore this error type because it isn't defined anywhere, and it's detecting a false positive
    - linters:
        - gosec
      text: "G304:"

    # Exclude lll issues for long lines with go:generate
    - linters:
        - lll
      source: "^//go:generate "

  # Independently from option ` + "`" + `exclude` + "`" + ` we use default exclude patterns,
  # it can be disabled by this option. To list all
  # excluded by default patterns execute ` + "`" + `golangci-lint run --help` + "`" + `.
  # Default value for this option is true.
  exclude-use-default: false

  # Maximum issues count per one linter. Set to 0 to disable. Default is 50.
  max-issues-per-linter: 0

  # Maximum count of issues with the same text. Set to 0 to disable. Default is 3.
  max-same-issues: 0

  # Show only new issues: if there are unstaged changes or untracked files,
  # only those changes are analyzed, else only changes in HEAD~ are analyzed.
  # It's a super-useful option for integration of golangci-lint into existing
  # large codebase. It's not practical to fix all existing issues at the moment
  # of integration: much better don't allow issues in new code.
  # Default is false.
  new: false

  #
  # # Show only new issues created after git revision ` + "`" + `REV` + "`" + `
  # new-from-rev: REV
  #
  # # Show only new issues created in git patch with set file path.
  # new-from-patch: path/to/patch/file
  #
`
		actual := golangCILintDotYAML(project)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
