package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_implementationTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := implementationTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	oauth2v3 "gopkg.in/oauth2.v3"
	errors1 "gopkg.in/oauth2.v3/errors"
	"net/http"
	"testing"
)

const (
	apiURLPrefix = "/api/v1"
)

func TestService_OAuth2InternalErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		expected := errors.New("blah")

		actual := s.OAuth2InternalErrorHandler(expected)
		assert.Equal(t, expected, actual.Error)
	})
}

func TestService_OAuth2ResponseErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &errors1.Response{}
		buildTestService(t).OAuth2ResponseErrorHandler(exampleInput)
	})
}

func TestService_AuthorizeScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t)
		res := httptest.NewRecorder()

		req = req.WithContext(
			context.WithValue(req.Context(), v1.OAuth2ClientKey, exampleOAuth2Client),
		)
		req.URL.Path = fmt.Sprintf("%s/%s", apiURLPrefix, exampleOAuth2Client.Scopes[0])
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, exampleOAuth2Client.Scopes[0], actual)
	})

	T.Run("without client attached to request but with client ID attached", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()

		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		req.URL.Path = fmt.Sprintf("%s/%s", apiURLPrefix, exampleOAuth2Client.Scopes[0])
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, exampleOAuth2Client.Scopes[0], actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request and now rows found fetching client info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v1.OAuth2Client)(nil), sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request and error fetching client info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request", func(t *testing.T) {
		s := buildTestService(t)
		req := buildRequest(t)
		res := httptest.NewRecorder()
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Empty(t, actual)
	})

	T.Run("with invalid scope & client ID but no client", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("%s/blah", apiURLPrefix)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_UserAuthorizationHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := fmt.Sprintf("%d", exampleOAuth2Client.BelongsToUser)

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), v1.OAuth2ClientKey, exampleOAuth2Client),
		)

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	T.Run("without client attached to request", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		expected := fmt.Sprintf("%d", exampleUser.ID)

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	T.Run("with no user info attached", func(t *testing.T) {
		s := buildTestService(t)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})
}

func TestService_ClientAuthorizedHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleGrant := oauth2v3.AuthorizationCode
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with password credentials grant", func(t *testing.T) {
		s := buildTestService(t)
		exampleGrant := oauth2v3.PasswordCredentials

		actual, err := s.ClientAuthorizedHandler("ID", exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)
		exampleGrant := oauth2v3.AuthorizationCode
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with disallowed implicit", func(t *testing.T) {
		s := buildTestService(t)

		exampleGrant := oauth2v3.Implicit
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_ClientScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleOAuth2Client.Scopes[0])
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleOAuth2Client.Scopes[0])
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without valid scope", func(t *testing.T) {
		s := buildTestService(t)

		exampleScope := "halb"
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleScope)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationTestConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildImplementationTestConstantDefs()

		expected := `
package example

import ()

const (
	apiURLPrefix = "/api/v1"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_OAuth2InternalErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestService_OAuth2InternalErrorHandler()

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestService_OAuth2InternalErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		expected := errors.New("blah")

		actual := s.OAuth2InternalErrorHandler(expected)
		assert.Equal(t, expected, actual.Error)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_OAuth2ResponseErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestService_OAuth2ResponseErrorHandler()

		expected := `
package example

import (
	errors "gopkg.in/oauth2.v3/errors"
	"testing"
)

func TestService_OAuth2ResponseErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &errors.Response{}
		buildTestService(t).OAuth2ResponseErrorHandler(exampleInput)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_AuthorizeScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_AuthorizeScopeHandler(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_AuthorizeScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t)
		res := httptest.NewRecorder()

		req = req.WithContext(
			context.WithValue(req.Context(), v1.OAuth2ClientKey, exampleOAuth2Client),
		)
		req.URL.Path = fmt.Sprintf("%s/%s", apiURLPrefix, exampleOAuth2Client.Scopes[0])
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, exampleOAuth2Client.Scopes[0], actual)
	})

	T.Run("without client attached to request but with client ID attached", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()

		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		req.URL.Path = fmt.Sprintf("%s/%s", apiURLPrefix, exampleOAuth2Client.Scopes[0])
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, exampleOAuth2Client.Scopes[0], actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request and now rows found fetching client info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v1.OAuth2Client)(nil), sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request and error fetching client info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without client attached to request", func(t *testing.T) {
		s := buildTestService(t)
		req := buildRequest(t)
		res := httptest.NewRecorder()
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Empty(t, actual)
	})

	T.Run("with invalid scope & client ID but no client", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("%s/blah", apiURLPrefix)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), clientIDKey, exampleOAuth2Client.ClientID),
		)
		actual, err := s.AuthorizeScopeHandler(res, req)

		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_UserAuthorizationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_UserAuthorizationHandler(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestService_UserAuthorizationHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := fmt.Sprintf("%d", exampleOAuth2Client.BelongsToUser)

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), v1.OAuth2ClientKey, exampleOAuth2Client),
		)

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	T.Run("without client attached to request", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		expected := fmt.Sprintf("%d", exampleUser.ID)

		req := buildRequest(t)
		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	T.Run("with no user info attached", func(t *testing.T) {
		s := buildTestService(t)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		actual, err := s.UserAuthorizationHandler(res, req)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_ClientAuthorizedHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_ClientAuthorizedHandler(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	oauth2v3 "gopkg.in/oauth2.v3"
	"testing"
)

func TestService_ClientAuthorizedHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleGrant := oauth2v3.AuthorizationCode
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with password credentials grant", func(t *testing.T) {
		s := buildTestService(t)
		exampleGrant := oauth2v3.PasswordCredentials

		actual, err := s.ClientAuthorizedHandler("ID", exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)
		exampleGrant := oauth2v3.AuthorizationCode
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with disallowed implicit", func(t *testing.T) {
		s := buildTestService(t)

		exampleGrant := oauth2v3.Implicit
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientAuthorizedHandler(stringID, exampleGrant)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_ClientScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_ClientScopeHandler(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestService_ClientScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleOAuth2Client.Scopes[0])
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleOAuth2Client.Scopes[0])
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without valid scope", func(t *testing.T) {
		s := buildTestService(t)

		exampleScope := "halb"
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		stringID := fmt.Sprintf("%d", exampleOAuth2Client.ID)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			stringID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		actual, err := s.ClientScopeHandler(stringID, exampleScope)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
