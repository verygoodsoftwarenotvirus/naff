package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func AddImports(proj *models.Project, file *jen.File, includeEmbedAnonymously bool) {
	pkgRoot := proj.OutputPath

	file.ImportAlias(proj.HTTPClientPackage(), "client")

	file.ImportAlias(proj.DatabasePackage(), "database")

	if includeEmbedAnonymously {
		file.Anon("embed")
	} else {
		file.ImportName("embed", "embed")
	}

	file.ImportAlias(proj.InternalAuthenticationPackage("mock"), "mockauth")
	file.ImportAlias(proj.EncodingPackage("mock"), "mockencoding")
	file.ImportAlias(proj.InternalRoutingPackage("mock"), "mockrouting")
	file.ImportAlias(proj.InternalSearchPackage("mock"), "mocksearch")
	file.ImportAlias(proj.MetricsPackage("mock"), "mockmetrics")

	file.ImportAlias(proj.DatabasePackage("config"), "dbconfig")
	file.ImportAlias(proj.DatabasePackage("client"), "dbclient")
	file.ImportAlias(proj.TypesPackage("mock"), "mocktypes")

	file.ImportAlias(filepath.Join(pkgRoot, "server"), "httpserver")

	file.ImportAlias(proj.APIClientsServicePackage(), "apiclientsservice")
	file.ImportAlias(proj.AccountsServicePackage(), "accountsservice")
	file.ImportAlias(proj.AdminServicePackage(), "adminservice")
	file.ImportAlias(proj.AuthServicePackage(), "authservice")
	file.ImportAlias(proj.FrontendServicePackage(), "frontendservice")
	file.ImportAlias(proj.UsersServicePackage(), "usersservice")
	file.ImportAlias(proj.WebhooksServicePackage(), "webhooksservice")
	file.ImportAlias(proj.WebsocketsServicePackage(), "websocketsservice")

	for _, typ := range proj.DataTypes {
		pn := typ.Name.PackageName()
		file.ImportAlias(filepath.Join(pkgRoot, "internal", "services", pn), fmt.Sprintf("%sservice", pn))
	}

	file.ImportAlias(proj.TestUtilsPackage(), "testutils")

	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/newsman/mock", "mocknewsman")

	file.ImportAlias("github.com/go-ozzo/ozzo-validation/v4", "validation")
	file.ImportAlias(constants.FakeLibrary, "fake")

	// databases
	file.ImportAlias("github.com/lib/pq", "postgres")
	file.ImportAlias("github.com/mattn/go-sqlite3", "sqlite")
	file.ImportName("github.com/go-sql-driver/mysql", "mysql")

	file.ImportNames(map[string]string{
		"context":                                "context",
		"fmt":                                    "fmt",
		"net/http":                               "http",
		"net/http/httputil":                      "httputil",
		"path":                                   "path",
		"path/filepath":                          "filepath",
		"errors":                                 "errors",
		"net/url":                                "url",
		"strings":                                "strings",
		"time":                                   "time",
		"bytes":                                  "bytes",
		"encoding/json":                          "json",
		"io":                                     "io",
		"io/ioutil":                              "ioutil",
		"reflect":                                "reflect",
		proj.InternalMessageQueueConfigPackage(): "msgconfig",
		proj.InternalMessageQueuePackage():       "messagequeue",
		proj.InternalAuthenticationPackage():     "authentication",
		proj.InternalAuthorizationPackage():      "authorization",
		proj.ConfigPackage():                     "config",
		proj.EncodingPackage():                   "encoding",
		proj.MetricsPackage():                    "metrics",
		proj.InternalTracingPackage():            "tracing",
		proj.InternalSearchPackage():             "search",
		proj.InternalSecretsPackage():            "secrets",
		proj.InternalEventsPackage():             "events",
		proj.InternalSearchPackage("elasticsearch"):        "elasticsearch",
		proj.DatabasePackage("queriers", "mysql"):          "mysql",
		proj.DatabasePackage("queriers", "postgres"):       "postgres",
		proj.TypesPackage():                                "types",
		proj.TypesPackage("fakes"):                         "fakes",
		proj.TestsPackage("frontend"):                      "frontend",
		proj.TestsPackage("integration"):                   "integration",
		proj.InternalLoggingPackage():                      "logging",
		proj.InternalLoggingPackage("zerolog"):             "zerolog",
		constants.RBACLibrary:                              "gorbac",
		constants.TracingAttributionLibrary:                "attribute",
		constants.AssertionLibrary:                         "assert",
		constants.MustAssertPkg:                            "require",
		constants.TestSuitePackage:                         "suite",
		constants.MockPkg:                                  "mock",
		constants.TracingLibrary:                           "trace",
		constants.SessionManagerLibrary:                    "scs",
		"github.com/boombuler/barcode/qr":                  "qr",
		"github.com/DATA-DOG/go-sqlmock":                   "sqlmock",
		"github.com/GuiaBolso/darwin":                      "darwin",
		"github.com/Masterminds/squirrel":                  "squirrel",
		"github.com/boombuler/barcode":                     "barcode",
		"github.com/nicksnyder/go-i18n/v2/i18n":            "i18n",
		"github.com/emicklei/hazana":                       "hazana",
		"github.com/go-chi/chi/v5":                         "chi",
		"github.com/blevesearch/bleve/v2/analysis/lang/en": "en",
		"github.com/blevesearch/bleve/v2/mapping":          "mapping",
		"github.com/blevesearch/bleve/v2/search/searcher":  "searcher",
		constants.SearchLibrary:                            "bleve",
		"github.com/go-chi/chi/v5/middleware":              "middleware",
		"github.com/go-chi/cors":                           "cors",
		constants.DependencyInjectionPkg:                   "wire",
		"github.com/gorilla/securecookie":                  "securecookie",
		"github.com/heptiolabs/healthcheck":                "healthcheck",
		"github.com/moul/http2curl":                        "http2curl",
		"github.com/pquerna/otp":                           "otp",
		"github.com/spf13/afero":                           "afero",
		"github.com/spf13/viper":                           "viper",
		"github.com/tebeka/selenium":                       "selenium",
		"golang.org/x/crypto":                              "crypto",
		constants.FlagParsingLibrary:                       "flag",
		"github.com/pquerna/otp/totp":                      "totp",
		"golang.org/x/oauth2/clientcredentials":            "clientcredentials",
	})

	file.Newline()
}

type importDef struct {
	name,
	path string
}

type importList []importDef

type importSet struct {
	stdlibImports   []importDef
	localLibImports []importDef
	externalImports []importDef
}

func (s *importSet) render() string {
	x := fmt.Sprintf("import(\n\t")

	for _, i := range s.stdlibImports {
		x += fmt.Sprintf("%s %q\n\t", i.name, i.path)
	}
	x += "\n\n\t"

	for _, i := range s.externalImports {
		x += fmt.Sprintf("%s %q\n\t", i.name, i.path)
	}
	x += "\n\n\t"

	for _, i := range s.localLibImports {
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
		stdlibImports:   []importDef{},
		localLibImports: []importDef{},
		externalImports: []importDef{},
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
			y := importDef{path: strings.ReplaceAll(x, `"`, "")}
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
	if err = ioutil.WriteFile(filepath, []byte(fart), 0644); err != nil {
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
		"embed",
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
		"io/fs",
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
