package utils

import (
	"bytes"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func AddImports(proj *models.Project, file *jen.File) {
	pkgRoot := proj.OutputPath

	file.ImportAlias(proj.HTTPClientV1Package(), "client")

	file.ImportAlias(proj.DatabaseV1Package(), "database")

	file.ImportName(proj.InternalAuthV1Package(), "auth")
	file.ImportAlias(proj.InternalAuthV1Package("mock"), "mockauth")
	file.ImportName(proj.InternalConfigV1Package(), "config")
	file.ImportName(proj.InternalEncodingV1Package(), "encoding")
	file.ImportAlias(proj.InternalEncodingV1Package("mock"), "mockencoding")
	file.ImportName(proj.InternalMetricsV1Package(), "metrics")
	file.ImportAlias(proj.InternalMetricsV1Package("mock"), "mockmetrics")
	file.ImportName(proj.InternalTracingV1Package(), "tracing")
	file.ImportName(proj.InternalSearchV1Package(), "search")
	file.ImportName(proj.InternalSearchV1Package("bleve"), "bleve")
	file.ImportName(proj.InternalSearchV1Package("mock"), "mocksearch")

	file.ImportAlias(proj.DatabaseV1Package("client"), "dbclient")
	file.ImportName(proj.DatabaseV1Package("queriers", "mariadb"), "mariadb")
	file.ImportName(proj.DatabaseV1Package("queriers", "postgres"), "postgres")
	file.ImportName(proj.DatabaseV1Package("queriers", "sqlite"), "sqlite")

	file.ImportAlias(proj.ModelsV1Package(), "models")
	file.ImportAlias(proj.ModelsV1Package("mock"), "mockmodels")
	file.ImportAlias(proj.ModelsV1Package("fake"), "fakemodels")

	file.ImportAlias(filepath.Join(pkgRoot, "server", "v1"), "server")
	file.ImportAlias(filepath.Join(pkgRoot, "server", "v1", "http"), "httpserver")

	file.ImportAlias(proj.ServiceV1AuthPackage(), "authservice")
	file.ImportAlias(proj.ServiceV1FrontendPackage(), "frontendservice")
	file.ImportAlias(proj.ServiceV1OAuth2ClientsPackage(), "oauth2clientsservice")
	file.ImportAlias(proj.ServiceV1UsersPackage(), "usersservice")
	file.ImportAlias(proj.ServiceV1WebhooksPackage(), "webhooksservice")

	for _, typ := range proj.DataTypes {
		pn := typ.Name.PackageName()
		file.ImportAlias(filepath.Join(pkgRoot, "services/v1", pn), fmt.Sprintf("%sservice", pn))
	}

	file.ImportName(filepath.Join(pkgRoot, "tests", "v1", "frontend"), "frontend")
	file.ImportName(filepath.Join(pkgRoot, "tests", "v1", "integration"), "integration")
	file.ImportName(filepath.Join(pkgRoot, "tests", "v1", "load"), "load")
	file.ImportName(proj.TestUtilV1Package(), "testutil")

	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "mocknewsman")

	file.ImportName(constants.CoreOAuth2Pkg, "oauth2")
	file.ImportName(constants.LoggingPkg, "logging")
	file.ImportName(constants.NoopLoggingPkg, "noop")
	file.ImportName(filepath.Join(constants.LoggingPkg, "zerolog"), "zerolog")
	file.ImportName(constants.AssertPkg, "assert")
	file.ImportName(constants.MustAssertPkg, "require")
	file.ImportName(constants.MockPkg, "mock")
	file.ImportAlias(constants.FakeLibrary, "fake")
	file.ImportName(constants.TracingLibrary, "trace")
	file.ImportName(constants.SessionManagerLibrary, "scs")

	file.ImportName("go.opencensus.io/stats", "stats")
	file.ImportName("go.opencensus.io/stats/view", "view")

	file.ImportAlias("gopkg.in/oauth2.v3", "oauth2")
	file.ImportAlias("gopkg.in/oauth2.v3/models", "oauth2models")
	file.ImportAlias("gopkg.in/oauth2.v3/errors", "oauth2errors")
	file.ImportAlias("gopkg.in/oauth2.v3/server", "oauth2server")
	file.ImportAlias("gopkg.in/oauth2.v3/store", "oauth2store")

	file.ImportName("golang.org/x/crypto/bcrypt", "bcrypt")

	// databases
	file.ImportAlias("github.com/lib/pq", "postgres")
	file.ImportAlias("github.com/mattn/go-sqlite3", "sqlite")
	file.ImportName("github.com/go-sql-driver/mysql", "mysql")

	file.ImportNames(map[string]string{
		"context":           "context",
		"fmt":               "fmt",
		"net/http":          "http",
		"net/http/httputil": "httputil",
		"errors":            "errors",
		"net/url":           "url",
		"path":              "path",
		"strings":           "strings",
		"time":              "time",
		"bytes":             "bytes",
		"encoding/json":     "json",
		"io":                "io",
		"io/ioutil":         "ioutil",
		"reflect":           "reflect",

		"gopkg.in/oauth2.v3/manage":                    "manage",
		"github.com/boombuler/barcode/qr":              "qr",
		"contrib.go.opencensus.io/exporter/jaeger":     "jaeger",
		"contrib.go.opencensus.io/exporter/prometheus": "prometheus",
		"contrib.go.opencensus.io/integrations/ocsql":  "ocsql",
		"github.com/DATA-DOG/go-sqlmock":               "sqlmock",
		"github.com/GuiaBolso/darwin":                  "darwin",
		"github.com/Masterminds/squirrel":              "squirrel",
		"github.com/boombuler/barcode":                 "barcode",
		"github.com/emicklei/hazana":                   "hazana",
		"github.com/go-chi/chi":                        "chi",
		"github.com/go-chi/chi/middleware":             "middleware",
		"github.com/go-chi/cors":                       "cors",
		"github.com/google/wire":                       "wire",
		"github.com/gorilla/securecookie":              "securecookie",
		"github.com/heptiolabs/healthcheck":            "healthcheck",
		"github.com/moul/http2curl":                    "http2curl",
		"github.com/pquerna/otp":                       "otp",
		"github.com/spf13/afero":                       "afero",
		"github.com/spf13/viper":                       "viper",
		"github.com/tebeka/selenium":                   "selenium",
		"gitlab.com/verygoodsoftwarenotvirus/newsman":  "newsman",
		"go.opencensus.io":                             "opencensus",
		"golang.org/x/crypto":                          "crypto",
		"go.opencensus.io/plugin/ochttp":               "ochttp",
		"github.com/spf13/pflag":                       "flag",
		"github.com/pquerna/otp/totp":                  "totp",
		"golang.org/x/oauth2/clientcredentials":        "clientcredentials",
	})

	file.Line()
}

type mport struct {
	name,
	path string
}

type importList []mport

type importSet struct {
	stdlibImports   []mport
	localLibImports []mport
	externalImports []mport
}

func (s *importSet) render() string {
	x := fmt.Sprintf("import(\n\t")

	for _, i := range s.stdlibImports {
		x += fmt.Sprintf("%s %q\n\t", i.name, i.path)
	}
	x += "\n\n\t"
	for _, i := range s.localLibImports {
		x += fmt.Sprintf("%s %q\n\t", i.name, i.path)
	}
	x += "\n\n\t"
	for _, i := range s.externalImports {
		x += fmt.Sprintf("%s %q\n\t", i.name, i.path)
	}

	x += "\n)\n\n"
	return x
}

// Len is the number of elements in the collection.
func (l importList) Len() int {
	return len(l)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (l importList) Less(i, j int) bool {
	return l[i].path < l[j].path
}

// Swap swaps the elements with indexes i and j.
func (l importList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l importList) toSet(pkgRoot string) *importSet {
	is := &importSet{
		stdlibImports:   []mport{},
		localLibImports: []mport{},
		externalImports: []mport{},
	}

	for _, imp := range l {
		if importIsStdLib(imp.path) {
			is.stdlibImports = append(is.stdlibImports, imp)
		} else if strings.HasPrefix(imp.path, pkgRoot) {
			is.localLibImports = append(is.localLibImports, imp)
		} else {
			is.externalImports = append(is.externalImports, imp)
		}
	}

	return is
}

func FindAndFixImportBlock(pkgRoot, filepath string) error {
	var allImports importList

	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	allLines := bytes.Split(fileBytes, []byte("\n"))
	var startLine, endLine int

	for i, l := range allLines {
		if len(l) > 0 {
			line := strings.TrimSpace(string(l))
			re := regexp.MustCompile(`import\s*\(`)
			if re.MatchString(line) {
				startLine = i
			} else if line == ")" {
				endLine = i
				break
			}
		}
	}

	if startLine == 0 || endLine == 0 {
		return nil
		//return errors.New("something done goofed")
	}

	for _, imp := range allLines[startLine+1 : endLine] {
		if x := strings.TrimSpace(string(imp)); x != "" {
			y := mport{path: strings.ReplaceAll(x, `"`, "")}
			if z := strings.Split(y.path, " "); len(z) > 1 {
				y.name = z[0]
				y.path = z[1]
			}
			allImports = append(allImports, y)
		}
	}

	is := allImports.toSet(pkgRoot).render()
	head := string(bytes.Join(allLines[:startLine], []byte("\n")))
	newImportBlock := strings.Join(strings.Split(is, "\n"), "\n")
	tail := string(bytes.Join(allLines[endLine+1:], []byte("\n")))

	fart := fmt.Sprintf("%s\n%s\n%s", head, newImportBlock, tail)

	os.Remove(filepath)
	if err := ioutil.WriteFile(filepath, []byte(fart), 0644); err != nil {
		return err
	}

	return nil
}

func importIsStdLib(imp string) bool {
	stdlibImports := []string{
		"archive/tar",
		"archive/zip",
		"bufio",
		"bytes",
		"compress/bzip2",
		"compress/flate",
		"compress/gzip",
		"compress/lzw",
		"compress/zlib",
		"container/heap",
		"container/list",
		"container/ring",
		"context",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/des",
		"crypto/dsa",
		"crypto/ecdsa",
		"crypto/ed25519",
		"crypto/ed25519/internal/edwards25519",
		"crypto/elliptic",
		"crypto/hmac",
		"crypto/internal/randutil",
		"crypto/internal/subtle",
		"crypto/md5",
		"crypto/rand",
		"crypto/rc4",
		"crypto/rsa",
		"crypto/sha1",
		"crypto/sha256",
		"crypto/sha512",
		"crypto/subtle",
		"crypto/tls",
		"crypto/x509",
		"crypto/x509/pkix",
		"database/sql",
		"database/sql/driver",
		"debug/dwarf",
		"debug/elf",
		"debug/gosym",
		"debug/macho",
		"debug/pe",
		"debug/plan9obj",
		"encoding",
		"encoding/ascii85",
		"encoding/asn1",
		"encoding/base32",
		"encoding/base64",
		"encoding/binary",
		"encoding/csv",
		"encoding/gob",
		"encoding/hex",
		"encoding/json",
		"encoding/pem",
		"encoding/xml",
		"errors",
		"expvar",
		"flag",
		"fmt",
		"go/ast",
		"go/build",
		"go/constant",
		"go/doc",
		"go/format",
		"go/importer",
		"go/internal/gccgoimporter",
		"go/internal/gcimporter",
		"go/internal/srcimporter",
		"go/parser",
		"go/printer",
		"go/scanner",
		"go/token",
		"go/types",
		"hash",
		"hash/adler32",
		"hash/crc32",
		"hash/crc64",
		"hash/fnv",
		"html",
		"html/template",
		"image",
		"image/color",
		"image/color/palette",
		"image/draw",
		"image/gif",
		"image/internal/imageutil",
		"image/jpeg",
		"image/png",
		"index/suffixarray",
		"internal/bytealg",
		"internal/cfg",
		"internal/cpu",
		"internal/fmtsort",
		"internal/goroot",
		"internal/goversion",
		"internal/lazyregexp",
		"internal/lazytemplate",
		"internal/nettrace",
		"internal/oserror",
		"internal/poll",
		"internal/race",
		"internal/reflectlite",
		"internal/singleflight",
		"internal/syscall/unix",
		"internal/testenv",
		"internal/testlog",
		"internal/trace",
		"internal/xcoff",
		"io",
		"io/ioutil",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/bits",
		"math/cmplx",
		"math/rand",
		"mime",
		"mime/multipart",
		"mime/quotedprintable",
		"net",
		"net/http",
		"net/http/cgi",
		"net/http/cookiejar",
		"net/http/fcgi",
		"net/http/httptest",
		"net/http/httptrace",
		"net/http/httputil",
		"net/http/internal",
		"net/http/pprof",
		"net/internal/socktest",
		"net/mail",
		"net/rpc",
		"net/rpc/jsonrpc",
		"net/smtp",
		"net/textproto",
		"net/url",
		"os",
		"os/exec",
		"os/signal",
		"os/signal/internal/pty",
		"os/user",
		"path",
		"path/filepath",
		"plugin",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/debug",
		"runtime/internal/atomic",
		"runtime/internal/math",
		"runtime/internal/sys",
		"runtime/pprof",
		"runtime/pprof/internal/profile",
		"runtime/race",
		"runtime/trace",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"testing",
		"testing/internal/testdeps",
		"testing/iotest",
		"testing/quick",
		"text/scanner",
		"text/tabwriter",
		"text/template",
		"text/template/parse",
		"time",
		"unicode",
		"unicode/utf16",
		"unicode/utf8",
		"unsafe",
	}

	for _, i := range stdlibImports {
		if i == strings.TrimSpace(strings.ToLower(imp)) {
			return true
		}
	}
	return false
}
