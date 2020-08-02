package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_userDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := userDotGo(proj)

		expected := `
package example

import (
	"encoding/base32"
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	totp "github.com/pquerna/otp/totp"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"time"
)

// BuildFakeUser builds a faked User.
func BuildFakeUser() *v1.User {
	return &v1.User{
		ID:       uint64(v5.Uint32()),
		Username: v5.Username(),
		// HashedPassword: "",
		// Salt:           []byte(fake.Word()),
		TwoFactorSecret:           base32.StdEncoding.EncodeToString([]byte(v5.Password(false, true, true, false, false, 32))),
		TwoFactorSecretVerifiedOn: func(i uint64) *uint64 { return &i }(uint64(uint32(v5.Date().Unix()))),
		IsAdmin:                   false,
		CreatedOn:                 uint64(uint32(v5.Date().Unix())),
	}
}

// BuildDatabaseCreationResponse builds a faked UserCreationResponse.
func BuildDatabaseCreationResponse(user *v1.User) *v1.UserCreationResponse {
	return &v1.UserCreationResponse{
		ID:                    user.ID,
		Username:              user.Username,
		TwoFactorSecret:       user.TwoFactorSecret,
		PasswordLastChangedOn: user.PasswordLastChangedOn,
		IsAdmin:               user.IsAdmin,
		CreatedOn:             user.CreatedOn,
		LastUpdatedOn:         user.LastUpdatedOn,
		ArchivedOn:            user.ArchivedOn,
	}
}

// BuildFakeUserList builds a faked UserList.
func BuildFakeUserList() *v1.UserList {
	exampleUser1 := BuildFakeUser()
	exampleUser2 := BuildFakeUser()
	exampleUser3 := BuildFakeUser()

	return &v1.UserList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Users: []v1.User{
			*exampleUser1,
			*exampleUser2,
			*exampleUser3,
		},
	}
}

// BuildFakeUserCreationInput builds a faked UserCreationInput.
func BuildFakeUserCreationInput() *v1.UserCreationInput {
	exampleUser := BuildFakeUser()
	return &v1.UserCreationInput{
		Username: exampleUser.Username,
		Password: v5.Password(true, true, true, true, true, 32),
	}
}

// BuildFakeUserCreationInputFromUser builds a faked UserCreationInput.
func BuildFakeUserCreationInputFromUser(user *v1.User) *v1.UserCreationInput {
	return &v1.UserCreationInput{
		Username: user.Username,
		Password: v5.Password(true, true, true, true, true, 32),
	}
}

// BuildFakeUserDatabaseCreationInputFromUser builds a faked UserDatabaseCreationInput.
func BuildFakeUserDatabaseCreationInputFromUser(user *v1.User) v1.UserDatabaseCreationInput {
	return v1.UserDatabaseCreationInput{
		Username:        user.Username,
		HashedPassword:  user.HashedPassword,
		TwoFactorSecret: user.TwoFactorSecret,
	}
}

// BuildFakeUserLoginInputFromUser builds a faked UserLoginInput.
func BuildFakeUserLoginInputFromUser(user *v1.User) *v1.UserLoginInput {
	return &v1.UserLoginInput{
		Username:  user.Username,
		Password:  v5.Password(true, true, true, true, true, 32),
		TOTPToken: fmt.Sprintf("0%s", v5.Zip()),
	}
}

// BuildFakePasswordUpdateInput builds a faked PasswordUpdateInput.
func BuildFakePasswordUpdateInput() *v1.PasswordUpdateInput {
	return &v1.PasswordUpdateInput{
		NewPassword:     v5.Password(true, true, true, true, true, 32),
		CurrentPassword: v5.Password(true, true, true, true, true, 32),
		TOTPToken:       fmt.Sprintf("0%s", v5.Zip()),
	}
}

// BuildFakeTOTPSecretRefreshInput builds a faked TOTPSecretRefreshInput.
func BuildFakeTOTPSecretRefreshInput() *v1.TOTPSecretRefreshInput {
	return &v1.TOTPSecretRefreshInput{
		CurrentPassword: v5.Password(true, true, true, true, true, 32),
		TOTPToken:       fmt.Sprintf("0%s", v5.Zip()),
	}
}

// BuildFakeTOTPSecretValidationInputForUser builds a faked TOTPSecretVerificationInput for a given user
func BuildFakeTOTPSecretValidationInputForUser(user *v1.User) *v1.TOTPSecretVerificationInput {
	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	if err != nil {
		log.Panicf("error generating TOTP token for fake user: %v", err)
	}

	return &v1.TOTPSecretVerificationInput{
		UserID:    user.ID,
		TOTPToken: token,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUser(proj)

		expected := `
package example

import (
	"encoding/base32"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUser builds a faked User.
func BuildFakeUser() *v1.User {
	return &v1.User{
		ID:       uint64(v5.Uint32()),
		Username: v5.Username(),
		// HashedPassword: "",
		// Salt:           []byte(fake.Word()),
		TwoFactorSecret:           base32.StdEncoding.EncodeToString([]byte(v5.Password(false, true, true, false, false, 32))),
		TwoFactorSecretVerifiedOn: func(i uint64) *uint64 { return &i }(uint64(uint32(v5.Date().Unix()))),
		IsAdmin:                   false,
		CreatedOn:                 uint64(uint32(v5.Date().Unix())),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildDatabaseCreationResponse(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildDatabaseCreationResponse(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildDatabaseCreationResponse builds a faked UserCreationResponse.
func BuildDatabaseCreationResponse(user *v1.User) *v1.UserCreationResponse {
	return &v1.UserCreationResponse{
		ID:                    user.ID,
		Username:              user.Username,
		TwoFactorSecret:       user.TwoFactorSecret,
		PasswordLastChangedOn: user.PasswordLastChangedOn,
		IsAdmin:               user.IsAdmin,
		CreatedOn:             user.CreatedOn,
		LastUpdatedOn:         user.LastUpdatedOn,
		ArchivedOn:            user.ArchivedOn,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUserList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUserList(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUserList builds a faked UserList.
func BuildFakeUserList() *v1.UserList {
	exampleUser1 := BuildFakeUser()
	exampleUser2 := BuildFakeUser()
	exampleUser3 := BuildFakeUser()

	return &v1.UserList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Users: []v1.User{
			*exampleUser1,
			*exampleUser2,
			*exampleUser3,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUserCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUserCreationInput(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUserCreationInput builds a faked UserCreationInput.
func BuildFakeUserCreationInput() *v1.UserCreationInput {
	exampleUser := BuildFakeUser()
	return &v1.UserCreationInput{
		Username: exampleUser.Username,
		Password: v5.Password(true, true, true, true, true, 32),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUserCreationInputFromUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUserCreationInputFromUser(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUserCreationInputFromUser builds a faked UserCreationInput.
func BuildFakeUserCreationInputFromUser(user *v1.User) *v1.UserCreationInput {
	return &v1.UserCreationInput{
		Username: user.Username,
		Password: v5.Password(true, true, true, true, true, 32),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUserDatabaseCreationInputFromUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUserDatabaseCreationInputFromUser(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUserDatabaseCreationInputFromUser builds a faked UserDatabaseCreationInput.
func BuildFakeUserDatabaseCreationInputFromUser(user *v1.User) v1.UserDatabaseCreationInput {
	return v1.UserDatabaseCreationInput{
		Username:        user.Username,
		HashedPassword:  user.HashedPassword,
		TwoFactorSecret: user.TwoFactorSecret,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeUserLoginInputFromUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeUserLoginInputFromUser(proj)

		expected := `
package example

import (
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeUserLoginInputFromUser builds a faked UserLoginInput.
func BuildFakeUserLoginInputFromUser(user *v1.User) *v1.UserLoginInput {
	return &v1.UserLoginInput{
		Username:  user.Username,
		Password:  v5.Password(true, true, true, true, true, 32),
		TOTPToken: fmt.Sprintf("0%s", v5.Zip()),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakePasswordUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakePasswordUpdateInput(proj)

		expected := `
package example

import (
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakePasswordUpdateInput builds a faked PasswordUpdateInput.
func BuildFakePasswordUpdateInput() *v1.PasswordUpdateInput {
	return &v1.PasswordUpdateInput{
		NewPassword:     v5.Password(true, true, true, true, true, 32),
		CurrentPassword: v5.Password(true, true, true, true, true, 32),
		TOTPToken:       fmt.Sprintf("0%s", v5.Zip()),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeTOTPSecretRefreshInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeTOTPSecretRefreshInput(proj)

		expected := `
package example

import (
	"fmt"
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeTOTPSecretRefreshInput builds a faked TOTPSecretRefreshInput.
func BuildFakeTOTPSecretRefreshInput() *v1.TOTPSecretRefreshInput {
	return &v1.TOTPSecretRefreshInput{
		CurrentPassword: v5.Password(true, true, true, true, true, 32),
		TOTPToken:       fmt.Sprintf("0%s", v5.Zip()),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeTOTPSecretValidationInputForUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFakeTOTPSecretValidationInputForUser(proj)

		expected := `
package example

import (
	totp "github.com/pquerna/otp/totp"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"time"
)

// BuildFakeTOTPSecretValidationInputForUser builds a faked TOTPSecretVerificationInput for a given user
func BuildFakeTOTPSecretValidationInputForUser(user *v1.User) *v1.TOTPSecretVerificationInput {
	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	if err != nil {
		log.Panicf("error generating TOTP token for fake user: %v", err)
	}

	return &v1.TOTPSecretVerificationInput{
		UserID:    user.ID,
		TOTPToken: token,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
