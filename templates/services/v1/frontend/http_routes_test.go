package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := httpRoutesDotGo(proj)

		expected := `
package example

import (
	"fmt"
	afero "github.com/spf13/afero"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

func (s *Service) buildStaticFileServer(fileDir string) (*afero.HttpFs, error) {
	var afs afero.Fs
	if s.config.CacheStaticFiles {
		afs = afero.NewMemMapFs()
		files, err := ioutil.ReadDir(fileDir)
		if err != nil {
			return nil, fmt.Errorf("reading directory for frontend files: %w", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			fp := filepath.Join(fileDir, file.Name())
			f, err := afs.Create(fp)
			if err != nil {
				return nil, fmt.Errorf("creating static file in memory: %w", err)
			}

			bs, err := ioutil.ReadFile(fp)
			if err != nil {
				return nil, fmt.Errorf("reading static file from directory: %w", err)
			}

			if _, err = f.Write(bs); err != nil {
				return nil, fmt.Errorf("loading static file into memory: %w", err)
			}

			if err = f.Close(); err != nil {
				s.logger.Error(err, "closing file while setting up static dir")
			}
		}
		afs = afero.NewReadOnlyFs(afs)
	} else {
		afs = afero.NewOsFs()
	}

	return afero.NewHttpFs(afs), nil
}

var (
	// Here is where you should put route regexes that need to be ignored by the static file server.
	// For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic.
	// information, such as ` + "`" + `/event/123` + "`" + `, you would want to put something like this below:
	// 		eventsFrontendPathRegex = regexp.MustCompile(` + "`" + `/event/\d+` + "`" + `)

	// itemsFrontendPathRegex matches URLs against our frontend router's specification for specific item routes.
	itemsFrontendPathRegex = regexp.MustCompile(` + "`" + `/items/\d+` + "`" + `)
)

// StaticDir builds a static directory handler.
func (s *Service) StaticDir(staticFilesDirectory string) (http.HandlerFunc, error) {
	fileDir, err := filepath.Abs(staticFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("determining absolute path of static files directory: %w", err)
	}

	httpFs, err := s.buildStaticFileServer(fileDir)
	if err != nil {
		return nil, fmt.Errorf("establishing static server filesystem: %w", err)
	}

	s.logger.WithValue("static_dir", fileDir).Debug("setting static file server")
	fs := http.StripPrefix("/", http.FileServer(httpFs.Dir(fileDir)))

	return func(res http.ResponseWriter, req *http.Request) {
		rl := s.logger.WithRequest(req)
		rl.Debug("static file requested")
		switch req.URL.Path {
		// list your frontend history routes here.
		case "/register",
			"/login",
			"/items",
			"/items/new",
			"/password/new":
			rl.Debug("rerouting")
			req.URL.Path = "/"
		}
		if itemsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting item request")
			req.URL.Path = "/"
		}

		fs.ServeHTTP(res, req)
	}, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFrontendBuildStaticFileServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildFrontendBuildStaticFileServer()

		expected := `
package example

import (
	"fmt"
	afero "github.com/spf13/afero"
	"io/ioutil"
	"path/filepath"
)

func (s *Service) buildStaticFileServer(fileDir string) (*afero.HttpFs, error) {
	var afs afero.Fs
	if s.config.CacheStaticFiles {
		afs = afero.NewMemMapFs()
		files, err := ioutil.ReadDir(fileDir)
		if err != nil {
			return nil, fmt.Errorf("reading directory for frontend files: %w", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			fp := filepath.Join(fileDir, file.Name())
			f, err := afs.Create(fp)
			if err != nil {
				return nil, fmt.Errorf("creating static file in memory: %w", err)
			}

			bs, err := ioutil.ReadFile(fp)
			if err != nil {
				return nil, fmt.Errorf("reading static file from directory: %w", err)
			}

			if _, err = f.Write(bs); err != nil {
				return nil, fmt.Errorf("loading static file into memory: %w", err)
			}

			if err = f.Close(); err != nil {
				s.logger.Error(err, "closing file while setting up static dir")
			}
		}
		afs = afero.NewReadOnlyFs(afs)
	} else {
		afs = afero.NewOsFs()
	}

	return afero.NewHttpFs(afs), nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFrontendVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFrontendVarDeclarations(proj)

		expected := `
package example

import (
	"regexp"
)

var (
	// Here is where you should put route regexes that need to be ignored by the static file server.
	// For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic.
	// information, such as ` + "`" + `/event/123` + "`" + `, you would want to put something like this below:
	// 		eventsFrontendPathRegex = regexp.MustCompile(` + "`" + `/event/\d+` + "`" + `)

	// itemsFrontendPathRegex matches URLs against our frontend router's specification for specific item routes.
	itemsFrontendPathRegex = regexp.MustCompile(` + "`" + `/items/\d+` + "`" + `)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFrontendStaticDir(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFrontendStaticDir(proj)

		expected := `
package example

import (
	"fmt"
	"net/http"
	"path/filepath"
)

// StaticDir builds a static directory handler.
func (s *Service) StaticDir(staticFilesDirectory string) (http.HandlerFunc, error) {
	fileDir, err := filepath.Abs(staticFilesDirectory)
	if err != nil {
		return nil, fmt.Errorf("determining absolute path of static files directory: %w", err)
	}

	httpFs, err := s.buildStaticFileServer(fileDir)
	if err != nil {
		return nil, fmt.Errorf("establishing static server filesystem: %w", err)
	}

	s.logger.WithValue("static_dir", fileDir).Debug("setting static file server")
	fs := http.StripPrefix("/", http.FileServer(httpFs.Dir(fileDir)))

	return func(res http.ResponseWriter, req *http.Request) {
		rl := s.logger.WithRequest(req)
		rl.Debug("static file requested")
		switch req.URL.Path {
		// list your frontend history routes here.
		case "/register",
			"/login",
			"/items",
			"/items/new",
			"/password/new":
			rl.Debug("rerouting")
			req.URL.Path = "/"
		}
		if itemsFrontendPathRegex.MatchString(req.URL.Path) {
			rl.Debug("rerouting item request")
			req.URL.Path = "/"
		}

		fs.ServeHTTP(res, req)
	}, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
