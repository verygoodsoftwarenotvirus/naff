package testutil

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_testutilDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := testutilDotGo(proj)

		//

		expected := `
package example

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	http2curl "github.com/moul/http2curl"
	totp "github.com/pquerna/otp/totp"
	http1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func init() {
	v5.Seed(time.Now().UnixNano())
}

// DetermineServiceURL returns the URL, if properly configured.
func DetermineServiceURL() string {
	ta := os.Getenv("TARGET_ADDRESS")
	if ta == "" {
		panic("must provide target address!")
	}

	u, err := url.Parse(ta)
	if err != nil {
		panic(err)
	}

	svcAddr := u.String()

	log.Printf("using target address: %q\n", svcAddr)
	return svcAddr
}

// EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other.
func EnsureServerIsUp(address string) {
	var (
		isDown           = true
		interval         = time.Second
		maxAttempts      = 50
		numberOfAttempts = 0
	)

	for isDown {
		if !IsUp(address) {
			log.Print("waiting before pinging again")
			time.Sleep(interval)
			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				log.Fatal("Maximum number of attempts made, something's gone awry")
			}
		} else {
			isDown = false
		}
	}
}

// IsUp can check if an instance of our server is alive.
func IsUp(address string) bool {
	uri := fmt.Sprintf("%s/_meta_/ready", address)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	if err = res.Body.Close(); err != nil {
		log.Println("error closing body")
	}

	return res.StatusCode == http.StatusOK
}

// CreateObligatoryUser creates a user for the sake of having an OAuth2 client.
func CreateObligatoryUser(address string, debug bool) (*v1.User, error) {
	ctx := context.Background()
	tu, parseErr := url.Parse(address)
	if parseErr != nil {
		return nil, parseErr
	}

	c, clientInitErr := http1.NewSimpleClient(ctx, tu, debug)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	// I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,
	// could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit
	username := v5.Username() + v5.HexColor() + v5.Country()
	in := &v1.UserCreationInput{
		Username: username,
		Password: v5.Password(true, true, true, true, true, 64),
	}

	ucr, userCreationErr := c.CreateUser(ctx, in)
	if userCreationErr != nil {
		return nil, userCreationErr
	} else if ucr == nil {
		return nil, errors.New("something happened")
	}

	token, tokenErr := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
	if tokenErr != nil {
		return nil, fmt.Errorf("generating totp code: %w", tokenErr)
	}

	if validationErr := c.VerifyTOTPSecret(ctx, ucr.ID, token); validationErr != nil {
		return nil, fmt.Errorf("verifying totp code: %w", validationErr)
	}

	u := &v1.User{
		ID:       ucr.ID,
		Username: ucr.Username,
		// this is a dirty trick to reuse most of this model,
		HashedPassword:        in.Password,
		TwoFactorSecret:       ucr.TwoFactorSecret,
		PasswordLastChangedOn: ucr.PasswordLastChangedOn,
		CreatedOn:             ucr.CreatedOn,
		LastUpdatedOn:         ucr.LastUpdatedOn,
		ArchivedOn:            ucr.ArchivedOn,
	}

	return u, nil
}

func buildURL(address string, parts ...string) string {
	tu, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		panic(err)
	}

	return tu.ResolveReference(u).String()
}

func getLoginCookie(serviceURL string, u *v1.User) (*http.Cookie, error) {
	uri := buildURL(serviceURL, "users", "login")
	code, err := totp.GenerateCode(strings.ToUpper(u.TwoFactorSecret), time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("generating totp token: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		uri,
		strings.NewReader(
			fmt.Sprintf(
				` + "`" + `
					{
						"username": %q,
						"password": %q,
						"totpToken": %q
					}
				` + "`" + `,
				u.Username,
				u.HashedPassword,
				code,
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		log.Println("error closing body")
	}

	cookies := res.Cookies()
	if len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, errors.New("no cookie found :(")
}

// CreateObligatoryClient creates the OAuth2 client we need for tests.
func CreateObligatoryClient(serviceURL string, u *v1.User) (*v1.OAuth2Client, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}

	firstOAuth2ClientURI := buildURL(serviceURL, "oauth2", "client")

	code, err := totp.GenerateCode(
		strings.ToUpper(u.TwoFactorSecret),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		firstOAuth2ClientURI,
		strings.NewReader(fmt.Sprintf(` + "`" + `
	{
		"username": %q,
		"password": %q,
		"totpToken": %q,
		"belongsToUser": %d,
		"scopes": ["*"]
	}
		` + "`" + `, u.Username, u.HashedPassword, code, u.ID)),
	)
	if err != nil {
		return nil, err
	}

	cookie, err := getLoginCookie(serviceURL, u)
	if err != nil || cookie == nil {
		log.Fatalf("\ncookie problems!\n\tcookie == nil: %v\n\t\t\t  err: %v\n\t", cookie == nil, err)
	}
	req.AddCookie(cookie)
	var o v1.OAuth2Client

	var command fmt.Stringer
	if command, err = http2curl.GetCurlCommand(req); err == nil {
		log.Println(command.String())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bad status: %d", res.StatusCode)
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	bdump, err := httputil.DumpResponse(res, true)
	if err == nil && req.Method != http.MethodGet {
		log.Println(string(bdump))
	}

	return &o, json.NewDecoder(res.Body).Decode(&o)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildInit()

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	"time"
)

func init() {
	v5.Seed(time.Now().UnixNano())
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDetermineServiceURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildDetermineServiceURL()

		expected := `
package example

import (
	"log"
	"net/url"
	"os"
)

// DetermineServiceURL returns the URL, if properly configured.
func DetermineServiceURL() string {
	ta := os.Getenv("TARGET_ADDRESS")
	if ta == "" {
		panic("must provide target address!")
	}

	u, err := url.Parse(ta)
	if err != nil {
		panic(err)
	}

	svcAddr := u.String()

	log.Printf("using target address: %q\n", svcAddr)
	return svcAddr
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEnsureServerIsUp(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildEnsureServerIsUp()

		expected := `
package example

import (
	"log"
	"time"
)

// EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other.
func EnsureServerIsUp(address string) {
	var (
		isDown           = true
		interval         = time.Second
		maxAttempts      = 50
		numberOfAttempts = 0
	)

	for isDown {
		if !IsUp(address) {
			log.Print("waiting before pinging again")
			time.Sleep(interval)
			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				log.Fatal("Maximum number of attempts made, something's gone awry")
			}
		} else {
			isDown = false
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIsUp(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildIsUp()

		expected := `
package example

import (
	"fmt"
	"log"
	"net/http"
)

// IsUp can check if an instance of our server is alive.
func IsUp(address string) bool {
	uri := fmt.Sprintf("%s/_meta_/ready", address)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	if err = res.Body.Close(); err != nil {
		log.Println("error closing body")
	}

	return res.StatusCode == http.StatusOK
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateObligatoryUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateObligatoryUser(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	totp "github.com/pquerna/otp/totp"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/url"
	"time"
)

// CreateObligatoryUser creates a user for the sake of having an OAuth2 client.
func CreateObligatoryUser(address string, debug bool) (*v1.User, error) {
	ctx := context.Background()
	tu, parseErr := url.Parse(address)
	if parseErr != nil {
		return nil, parseErr
	}

	c, clientInitErr := http.NewSimpleClient(ctx, tu, debug)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	// I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,
	// could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit
	username := v5.Username() + v5.HexColor() + v5.Country()
	in := &v1.UserCreationInput{
		Username: username,
		Password: v5.Password(true, true, true, true, true, 64),
	}

	ucr, userCreationErr := c.CreateUser(ctx, in)
	if userCreationErr != nil {
		return nil, userCreationErr
	} else if ucr == nil {
		return nil, errors.New("something happened")
	}

	token, tokenErr := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
	if tokenErr != nil {
		return nil, fmt.Errorf("generating totp code: %w", tokenErr)
	}

	if validationErr := c.VerifyTOTPSecret(ctx, ucr.ID, token); validationErr != nil {
		return nil, fmt.Errorf("verifying totp code: %w", validationErr)
	}

	u := &v1.User{
		ID:       ucr.ID,
		Username: ucr.Username,
		// this is a dirty trick to reuse most of this model,
		HashedPassword:        in.Password,
		TwoFactorSecret:       ucr.TwoFactorSecret,
		PasswordLastChangedOn: ucr.PasswordLastChangedOn,
		CreatedOn:             ucr.CreatedOn,
		LastUpdatedOn:         ucr.LastUpdatedOn,
		ArchivedOn:            ucr.ArchivedOn,
	}

	return u, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildURL()

		expected := `
package example

import (
	"net/url"
	"strings"
)

func buildURL(address string, parts ...string) string {
	tu, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		panic(err)
	}

	return tu.ResolveReference(u).String()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetLoginCookie(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetLoginCookie(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	totp "github.com/pquerna/otp/totp"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"net/http"
	"strings"
	"time"
)

func getLoginCookie(serviceURL string, u *v1.User) (*http.Cookie, error) {
	uri := buildURL(serviceURL, "users", "login")
	code, err := totp.GenerateCode(strings.ToUpper(u.TwoFactorSecret), time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("generating totp token: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		uri,
		strings.NewReader(
			fmt.Sprintf(
				` + "`" + `
					{
						"username": %q,
						"password": %q,
						"totpToken": %q
					}
				` + "`" + `,
				u.Username,
				u.HashedPassword,
				code,
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		log.Println("error closing body")
	}

	cookies := res.Cookies()
	if len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, errors.New("no cookie found :(")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateObligatoryClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateObligatoryClient(proj)

		expected := `
package example

import (
	"encoding/json"
	"errors"
	"fmt"
	http2curl "github.com/moul/http2curl"
	totp "github.com/pquerna/otp/totp"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"net/http"
	"strings"
	"time"
)

// CreateObligatoryClient creates the OAuth2 client we need for tests.
func CreateObligatoryClient(serviceURL string, u *v1.User) (*v1.OAuth2Client, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}

	firstOAuth2ClientURI := buildURL(serviceURL, "oauth2", "client")

	code, err := totp.GenerateCode(
		strings.ToUpper(u.TwoFactorSecret),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		firstOAuth2ClientURI,
		strings.NewReader(fmt.Sprintf(` + "`" + `
	{
		"username": %q,
		"password": %q,
		"totpToken": %q,
		"belongsToUser": %d,
		"scopes": ["*"]
	}
		` + "`" + `, u.Username, u.HashedPassword, code, u.ID)),
	)
	if err != nil {
		return nil, err
	}

	cookie, err := getLoginCookie(serviceURL, u)
	if err != nil || cookie == nil {
		log.Fatalf("\ncookie problems!\n\tcookie == nil: %v\n\t\t\t  err: %v\n\t", cookie == nil, err)
	}
	req.AddCookie(cookie)
	var o v1.OAuth2Client

	var command fmt.Stringer
	if command, err = http2curl.GetCurlCommand(req); err == nil {
		log.Println(command.String())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bad status: %d", res.StatusCode)
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	bdump, err := httputil.DumpResponse(res, true)
	if err == nil && req.Method != http.MethodGet {
		log.Println(string(bdump))
	}

	return &o, json.NewDecoder(res.Body).Decode(&o)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
