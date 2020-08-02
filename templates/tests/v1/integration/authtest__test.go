package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := authTestDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func loginUser(ctx context.Context, t *testing.T, username, password, totpSecret string) *http.Cookie {
	loginURL := fmt.Sprintf("%s://%s:%s/users/login", todoClient.URL.Scheme, todoClient.URL.Hostname(), todoClient.URL.Port())

	code, err := totp.GenerateCode(strings.ToUpper(totpSecret), time.Now().UTC())
	assert.NoError(t, err)

	bodyStr := fmt.Sprintf(` + "`" + `
	{
		"username": %q,
		"password": %q,
		"totpToken": %q
	}
` + "`" + `, username, password, code)

	body := strings.NewReader(bodyStr)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "login should be successful")

	cookies := resp.Cookies()
	if len(cookies) == 1 {
		return cookies[0]
	}
	t.Logf("wrong number of cookies found: %d", len(cookies))
	t.FailNow()

	return nil
}

func TestAuth(test *testing.T) {
	test.Run("should be able to login", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		secretVerificationToken, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, ucr.ID, secretVerificationToken))

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 1)
	})

	test.Run("should be able to logout", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		secretVerificationToken, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, ucr.ID, secretVerificationToken))

		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		// extract cookie.
		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		loginCookie := cookies[0]

		// build logout request.
		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(loginCookie)

		// execute logout request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	test.Run("login request without body fails", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	test.Run("should not be able to log in with the wrong password", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		// create login request.
		var badPassword string
		for _, v := range exampleUserCreationInput.Password {
			badPassword = string(v) + badPassword
		}

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  badPassword,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	test.Run("should not be able to login as someone that doesn't exist", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		s, err := randString()
		require.NoError(t, err)

		token, err := totp.GenerateCode(s, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 0)
	})

	test.Run("should not be able to login without validating TOTP secret", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 0)
	})

	test.Run("should reject an unauthenticated request", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, todoClient.BuildURL(nil, "webhooks"), nil)
		assert.NoError(t, err)

		res, err := todoClient.PlainClient().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	test.Run("should be able to change password", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		user, ui, cookie := buildDummyUser(ctx, test)
		require.NotNil(test, cookie)

		// create login request.
		var backwardsPass string
		for _, v := range ui.Password {
			backwardsPass = string(v) + backwardsPass
		}

		// create password update request.
		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.PasswordUpdateInput{
			CurrentPassword: ui.Password,
			TOTPToken:       token,
			NewPassword:     backwardsPass,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/password/new"

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), body)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		// execute password update request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode)

		// logout.

		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// create login request.
		newToken, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, newToken, err)
		l, err := json.Marshal(&v1.UserLoginInput{
			Username:  user.Username,
			Password:  backwardsPass,
			TOTPToken: newToken,
		})
		require.NoError(t, err)
		body = bytes.NewReader(l)

		u3, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u3.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u3.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		assert.NotEqual(t, cookie, cookies[0])
	})

	test.Run("should be able to validate a 2FA token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		userInput := fake.BuildFakeUserCreationInput()
		user, err := todoClient.CreateUser(ctx, userInput)
		assert.NotNil(t, user)
		require.NoError(t, err)

		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, token))
	})

	test.Run("should reject attempt to validate an invalid 2FA token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		userInput := fake.BuildFakeUserCreationInput()
		user, err := todoClient.CreateUser(ctx, userInput)
		assert.NotNil(t, user)
		require.NoError(t, err)

		assert.Error(t, todoClient.VerifyTOTPSecret(ctx, user.ID, "NOTREAL"))
	})

	test.Run("should be able to change 2FA Token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		user, ui, cookie := buildDummyUser(ctx, test)
		require.NotNil(test, cookie)

		// create TOTP secret update request.
		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		ir := &v1.TOTPSecretRefreshInput{
			CurrentPassword: ui.Password,
			TOTPToken:       token,
		}
		out, err := json.Marshal(ir)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/totp_secret/new"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		// execute TOTP secret update request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode)

		// load user response.
		r := &v1.TOTPSecretRefreshResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(r))
		require.NotEqual(t, user.TwoFactorSecret, r.TwoFactorSecret)

		secretVerificationToken, err := totp.GenerateCode(r.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, secretVerificationToken))

		// logout.

		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// create login request.
		newToken, err := totp.GenerateCode(r.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, newToken, err)
		l, err := json.Marshal(&v1.UserLoginInput{
			Username:  user.Username,
			Password:  ui.Password,
			TOTPToken: newToken,
		})
		require.NoError(t, err)
		body = bytes.NewReader(l)

		u3, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u3.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u3.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		assert.NotEqual(t, cookie, cookies[0])
	})

	test.Run("should accept a login cookie if a token is missing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		_, _, cookie := buildDummyUser(ctx, test)
		assert.NotNil(t, cookie)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, todoClient.BuildURL(nil, "webhooks"), nil)
		assert.NoError(t, err)
		req.AddCookie(cookie)

		res, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	test.Run("should only allow users to see their own content", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user and oauth2 client A.
		userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
		require.NoError(t, err)

		ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
		require.NoError(t, err)

		clientA, err := http1.NewClient(
			ctx,
			ca.ClientID,
			ca.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			ca.Scopes,
			true,
		)
		checkValueAndError(test, clientA, err)

		// create webhook for user A.
		wciA := fake.BuildFakeWebhookCreationInput()
		webhookA, err := clientA.CreateWebhook(ctx, wciA)
		checkValueAndError(t, webhookA, err)

		// create user and oauth2 client B.
		userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
		require.NoError(t, err)

		cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
		require.NoError(t, err)

		clientB, err := http1.NewClient(
			ctx,
			cb.ClientID,
			cb.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			cb.Scopes,
			true,
		)
		checkValueAndError(test, clientB, err)

		// create webhook for user B.
		wciB := fake.BuildFakeWebhookCreationInput()
		webhookB, err := clientB.CreateWebhook(ctx, wciB)
		checkValueAndError(t, webhookB, err)

		i, err := clientB.GetWebhook(ctx, webhookA.ID)
		assert.Nil(t, i)
		assert.Error(t, err, "should experience error trying to fetch entry they're not authorized for")

		// Clean up.
		assert.NoError(t, todoClient.ArchiveWebhook(ctx, webhookA.ID))
		assert.NoError(t, todoClient.ArchiveWebhook(ctx, webhookB.ID))
	})

	test.Run("should only allow clients with a given scope to see that scope's content", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		x, y, cookie := buildDummyUser(ctx, test)
		assert.NotNil(t, cookie)

		input := buildDummyOAuth2ClientInput(test, x.Username, y.Password, x.TwoFactorSecret)
		input.Scopes = []string{"absolutelynevergonnaexistascopelikethis"}
		premade, err := todoClient.CreateOAuth2Client(ctx, cookie, input)
		checkValueAndError(test, premade, err)

		c, err := http1.NewClient(
			ctx,
			premade.ClientID,
			premade.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			premade.Scopes,
			true,
		)
		checkValueAndError(test, c, err)

		i, err := c.GetOAuth2Clients(ctx, nil)
		assert.Nil(t, i)
		assert.Error(t, err, "should experience error trying to fetch entry they're not authorized for")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLoginUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildLoginUser(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func loginUser(ctx context.Context, t *testing.T, username, password, totpSecret string) *http.Cookie {
	loginURL := fmt.Sprintf("%s://%s:%s/users/login", todoClient.URL.Scheme, todoClient.URL.Hostname(), todoClient.URL.Port())

	code, err := totp.GenerateCode(strings.ToUpper(totpSecret), time.Now().UTC())
	assert.NoError(t, err)

	bodyStr := fmt.Sprintf(` + "`" + `
	{
		"username": %q,
		"password": %q,
		"totpToken": %q
	}
` + "`" + `, username, password, code)

	body := strings.NewReader(bodyStr)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "login should be successful")

	cookies := resp.Cookies()
	if len(cookies) == 1 {
		return cookies[0]
	}
	t.Logf("wrong number of cookies found: %d", len(cookies))
	t.FailNow()

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestAuth(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestAuth(proj)

		expected := `
package example

import (
	"bytes"
	"context"
	"encoding/json"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestAuth(test *testing.T) {
	test.Run("should be able to login", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		secretVerificationToken, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, ucr.ID, secretVerificationToken))

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 1)
	})

	test.Run("should be able to logout", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		secretVerificationToken, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, ucr.ID, secretVerificationToken))

		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		// extract cookie.
		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		loginCookie := cookies[0]

		// build logout request.
		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(loginCookie)

		// execute logout request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	test.Run("login request without body fails", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	test.Run("should not be able to log in with the wrong password", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		// create login request.
		var badPassword string
		for _, v := range exampleUserCreationInput.Password {
			badPassword = string(v) + badPassword
		}

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  badPassword,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	test.Run("should not be able to login as someone that doesn't exist", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		s, err := randString()
		require.NoError(t, err)

		token, err := totp.GenerateCode(s, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 0)
	})

	test.Run("should not be able to login without validating TOTP secret", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create a user.
		exampleUser := fake.BuildFakeUser()
		exampleUserCreationInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		req, err := todoClient.BuildCreateUserRequest(ctx, exampleUserCreationInput)
		checkValueAndError(t, req, err)

		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)

		// load user response.
		ucr := &v1.UserCreationResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(ucr))

		// create login request.
		token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.UserLoginInput{
			Username:  exampleUserCreationInput.Username,
			Password:  exampleUserCreationInput.Password,
			TOTPToken: token,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		cookies := res.Cookies()
		assert.Len(t, cookies, 0)
	})

	test.Run("should reject an unauthenticated request", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, todoClient.BuildURL(nil, "webhooks"), nil)
		assert.NoError(t, err)

		res, err := todoClient.PlainClient().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	test.Run("should be able to change password", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		user, ui, cookie := buildDummyUser(ctx, test)
		require.NotNil(test, cookie)

		// create login request.
		var backwardsPass string
		for _, v := range ui.Password {
			backwardsPass = string(v) + backwardsPass
		}

		// create password update request.
		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		r := &v1.PasswordUpdateInput{
			CurrentPassword: ui.Password,
			TOTPToken:       token,
			NewPassword:     backwardsPass,
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/password/new"

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), body)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		// execute password update request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode)

		// logout.

		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// create login request.
		newToken, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, newToken, err)
		l, err := json.Marshal(&v1.UserLoginInput{
			Username:  user.Username,
			Password:  backwardsPass,
			TOTPToken: newToken,
		})
		require.NoError(t, err)
		body = bytes.NewReader(l)

		u3, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u3.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u3.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		assert.NotEqual(t, cookie, cookies[0])
	})

	test.Run("should be able to validate a 2FA token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		userInput := fake.BuildFakeUserCreationInput()
		user, err := todoClient.CreateUser(ctx, userInput)
		assert.NotNil(t, user)
		require.NoError(t, err)

		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, token))
	})

	test.Run("should reject attempt to validate an invalid 2FA token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		userInput := fake.BuildFakeUserCreationInput()
		user, err := todoClient.CreateUser(ctx, userInput)
		assert.NotNil(t, user)
		require.NoError(t, err)

		assert.Error(t, todoClient.VerifyTOTPSecret(ctx, user.ID, "NOTREAL"))
	})

	test.Run("should be able to change 2FA Token", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		user, ui, cookie := buildDummyUser(ctx, test)
		require.NotNil(test, cookie)

		// create TOTP secret update request.
		token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, token, err)
		ir := &v1.TOTPSecretRefreshInput{
			CurrentPassword: ui.Password,
			TOTPToken:       token,
		}
		out, err := json.Marshal(ir)
		require.NoError(t, err)
		body := bytes.NewReader(out)

		u, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u.Path = "/users/totp_secret/new"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		// execute TOTP secret update request.
		res, err := todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode)

		// load user response.
		r := &v1.TOTPSecretRefreshResponse{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(r))
		require.NotEqual(t, user.TwoFactorSecret, r.TwoFactorSecret)

		secretVerificationToken, err := totp.GenerateCode(r.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, secretVerificationToken, err)

		assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, secretVerificationToken))

		// logout.

		u2, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u2.Path = "/users/logout"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u2.String(), nil)
		checkValueAndError(t, req, err)
		req.AddCookie(cookie)

		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// create login request.
		newToken, err := totp.GenerateCode(r.TwoFactorSecret, time.Now().UTC())
		checkValueAndError(t, newToken, err)
		l, err := json.Marshal(&v1.UserLoginInput{
			Username:  user.Username,
			Password:  ui.Password,
			TOTPToken: newToken,
		})
		require.NoError(t, err)
		body = bytes.NewReader(l)

		u3, err := url.Parse(todoClient.BuildURL(nil))
		require.NoError(t, err)
		u3.Path = "/users/login"

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u3.String(), body)
		checkValueAndError(t, req, err)

		// execute login request.
		res, err = todoClient.PlainClient().Do(req)
		checkValueAndError(t, res, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		cookies := res.Cookies()
		require.Len(t, cookies, 1)
		assert.NotEqual(t, cookie, cookies[0])
	})

	test.Run("should accept a login cookie if a token is missing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		_, _, cookie := buildDummyUser(ctx, test)
		assert.NotNil(t, cookie)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, todoClient.BuildURL(nil, "webhooks"), nil)
		assert.NoError(t, err)
		req.AddCookie(cookie)

		res, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	test.Run("should only allow users to see their own content", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user and oauth2 client A.
		userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
		require.NoError(t, err)

		ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
		require.NoError(t, err)

		clientA, err := http1.NewClient(
			ctx,
			ca.ClientID,
			ca.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			ca.Scopes,
			true,
		)
		checkValueAndError(test, clientA, err)

		// create webhook for user A.
		wciA := fake.BuildFakeWebhookCreationInput()
		webhookA, err := clientA.CreateWebhook(ctx, wciA)
		checkValueAndError(t, webhookA, err)

		// create user and oauth2 client B.
		userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
		require.NoError(t, err)

		cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
		require.NoError(t, err)

		clientB, err := http1.NewClient(
			ctx,
			cb.ClientID,
			cb.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			cb.Scopes,
			true,
		)
		checkValueAndError(test, clientB, err)

		// create webhook for user B.
		wciB := fake.BuildFakeWebhookCreationInput()
		webhookB, err := clientB.CreateWebhook(ctx, wciB)
		checkValueAndError(t, webhookB, err)

		i, err := clientB.GetWebhook(ctx, webhookA.ID)
		assert.Nil(t, i)
		assert.Error(t, err, "should experience error trying to fetch entry they're not authorized for")

		// Clean up.
		assert.NoError(t, todoClient.ArchiveWebhook(ctx, webhookA.ID))
		assert.NoError(t, todoClient.ArchiveWebhook(ctx, webhookB.ID))
	})

	test.Run("should only allow clients with a given scope to see that scope's content", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// create user.
		x, y, cookie := buildDummyUser(ctx, test)
		assert.NotNil(t, cookie)

		input := buildDummyOAuth2ClientInput(test, x.Username, y.Password, x.TwoFactorSecret)
		input.Scopes = []string{"absolutelynevergonnaexistascopelikethis"}
		premade, err := todoClient.CreateOAuth2Client(ctx, cookie, input)
		checkValueAndError(test, premade, err)

		c, err := http1.NewClient(
			ctx,
			premade.ClientID,
			premade.ClientSecret,
			todoClient.URL,
			noop.ProvideNoopLogger(),
			buildHTTPClient(),
			premade.Scopes,
			true,
		)
		checkValueAndError(test, c, err)

		i, err := c.GetOAuth2Clients(ctx, nil)
		assert.Nil(t, i)
		assert.Error(t, err, "should experience error trying to fetch entry they're not authorized for")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
