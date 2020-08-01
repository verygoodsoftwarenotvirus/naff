package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := oauth2ClientDotGo(proj)

		expected := `
package example

import (
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeOAuth2Client builds a faked OAuth2Client.
func BuildFakeOAuth2Client() *v1.OAuth2Client {
	return &v1.OAuth2Client{
		ID:           v5.Uint64(),
		Name:         v5.Word(),
		ClientID:     v5.UUID(),
		ClientSecret: v5.UUID(),
		RedirectURI:  v5.URL(),
		Scopes: []string{
			v5.Word(),
			v5.Word(),
			v5.Word(),
		},
		ImplicitAllowed: false,
		BelongsToUser:   v5.Uint64(),
		CreatedOn:       uint64(uint32(v5.Date().Unix())),
	}
}

// BuildFakeOAuth2ClientList builds a faked OAuth2ClientList.
func BuildFakeOAuth2ClientList() *v1.OAuth2ClientList {
	exampleOAuth2Client1 := BuildFakeOAuth2Client()
	exampleOAuth2Client2 := BuildFakeOAuth2Client()
	exampleOAuth2Client3 := BuildFakeOAuth2Client()

	return &v1.OAuth2ClientList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Clients: []v1.OAuth2Client{
			*exampleOAuth2Client1,
			*exampleOAuth2Client2,
			*exampleOAuth2Client3,
		},
	}
}

// BuildFakeOAuth2ClientCreationInputFromClient builds a faked OAuth2ClientCreationInput.
func BuildFakeOAuth2ClientCreationInputFromClient(client *v1.OAuth2Client) *v1.OAuth2ClientCreationInput {
	return &v1.OAuth2ClientCreationInput{
		UserLoginInput: v1.UserLoginInput{
			Username:  v5.Username(),
			Password:  v5.Password(true, true, true, true, true, 32),
			TOTPToken: fmt.Sprintf("0%s", v5.Zip()),
		},
		Name:          client.Name,
		Scopes:        client.Scopes,
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURI:   client.RedirectURI,
		BelongsToUser: client.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeOAuth2Client(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeOAuth2Client builds a faked OAuth2Client.
func BuildFakeOAuth2Client() *v1.OAuth2Client {
	return &v1.OAuth2Client{
		ID:           v5.Uint64(),
		Name:         v5.Word(),
		ClientID:     v5.UUID(),
		ClientSecret: v5.UUID(),
		RedirectURI:  v5.URL(),
		Scopes: []string{
			v5.Word(),
			v5.Word(),
			v5.Word(),
		},
		ImplicitAllowed: false,
		BelongsToUser:   v5.Uint64(),
		CreatedOn:       uint64(uint32(v5.Date().Unix())),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeOAuth2ClientList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeOAuth2ClientList(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeOAuth2ClientList builds a faked OAuth2ClientList.
func BuildFakeOAuth2ClientList() *v1.OAuth2ClientList {
	exampleOAuth2Client1 := BuildFakeOAuth2Client()
	exampleOAuth2Client2 := BuildFakeOAuth2Client()
	exampleOAuth2Client3 := BuildFakeOAuth2Client()

	return &v1.OAuth2ClientList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Clients: []v1.OAuth2Client{
			*exampleOAuth2Client1,
			*exampleOAuth2Client2,
			*exampleOAuth2Client3,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeOAuth2ClientCreationInputFromClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeOAuth2ClientCreationInputFromClient(proj)

		expected := `
package example

import (
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeOAuth2ClientCreationInputFromClient builds a faked OAuth2ClientCreationInput.
func BuildFakeOAuth2ClientCreationInputFromClient(client *v1.OAuth2Client) *v1.OAuth2ClientCreationInput {
	return &v1.OAuth2ClientCreationInput{
		UserLoginInput: v1.UserLoginInput{
			Username:  v5.Username(),
			Password:  v5.Password(true, true, true, true, true, 32),
			TOTPToken: fmt.Sprintf("0%s", v5.Zip()),
		},
		Name:          client.Name,
		Scopes:        client.Scopes,
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		RedirectURI:   client.RedirectURI,
		BelongsToUser: client.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
