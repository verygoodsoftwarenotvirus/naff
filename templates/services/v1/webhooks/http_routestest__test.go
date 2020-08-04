package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := httpRoutesTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_List(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhookList := fake.BuildFakeWebhookList()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return(exampleWebhookList, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.WebhookList)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with error fetching webhooks from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.WebhookList)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhookList := fake.BuildFakeWebhookList()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return(exampleWebhookList, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}

func TestValidateWebhook(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		assert.NoError(t, validateWebhook(exampleInput))
	})

	T.Run("with invalid method", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.Method = " MEATLOAF "

		assert.Error(t, validateWebhook(exampleInput))
	})

	T.Run("with invalid url", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.URL = "%zzzzz"

		assert.Error(t, validateWebhook(exampleInput))
	})
}

func TestWebhooksService_Create(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		mc := &mock3.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd, ed)
	})

	T.Run("with invalid webhook request", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleWebhook.URL = "%zzzzz"
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		mc := &mock3.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd, ed)
	})
}

func TestWebhooksService_Read(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with no such webhook in database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching webhook from database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}

func TestWebhooksService_Update(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error updating webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}

func TestWebhooksService_Archive(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		mc := &mock3.UnitCounter{}
		mc.On("Decrement", mock1.Anything).Return()
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(nil)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNoContent, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd)
	})

	T.Run("with no webhook in database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooksService_List(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooksService_List(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_List(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhookList := fake.BuildFakeWebhookList()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return(exampleWebhookList, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.WebhookList)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with error fetching webhooks from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.WebhookList)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhookList := fake.BuildFakeWebhookList()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhooks",
			mock1.Anything,
			exampleUser.ID,
			mock1.AnythingOfType("*models.QueryFilter"),
		).Return(exampleWebhookList, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.WebhookList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestValidateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestValidateWebhook(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestValidateWebhook(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		assert.NoError(t, validateWebhook(exampleInput))
	})

	T.Run("with invalid method", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.Method = " MEATLOAF "

		assert.Error(t, validateWebhook(exampleInput))
	})

	T.Run("with invalid url", func(t *testing.T) {
		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
		exampleInput.URL = "%zzzzz"

		assert.Error(t, validateWebhook(exampleInput))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooksService_Create(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooksService_Create(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_Create(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		mc := &mock.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd, ed)
	})

	T.Run("with invalid webhook request", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleWebhook.URL = "%zzzzz"
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		mc := &mock.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"CreateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.WebhookCreationInput"),
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooksService_Read(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooksService_Read(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_Read(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("with no such webhook in database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching webhook from database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooksService_Update(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooksService_Update(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_Update(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return((*v1.Webhook)(nil), errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error updating webhook", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeWebhookUpdateInputFromWebhook(exampleWebhook)

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock.WebhookDataManager{}
		wd.On(
			"GetWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(exampleWebhook, nil)

		wd.On(
			"UpdateWebhook",
			mock1.Anything,
			mock1.AnythingOfType("*models.Webhook"),
		).Return(nil)
		s.webhookDataManager = wd

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Webhook")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, wd, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooksService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooksService_Archive(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestWebhooksService_Archive(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		mc := &mock.UnitCounter{}
		mc.On("Decrement", mock1.Anything).Return()
		s.webhookCounter = mc

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(nil)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNoContent, res.Code)

		mock1.AssertExpectationsForObjects(t, mc, wd)
	})

	T.Run("with no webhook in database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(sql.ErrNoRows)
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleWebhook.BelongsToUser = exampleUser.ID

		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}

		s.webhookIDFetcher = func(req *http.Request) uint64 {
			return exampleWebhook.ID
		}

		wd := &mock2.WebhookDataManager{}
		wd.On(
			"ArchiveWebhook",
			mock1.Anything,
			exampleWebhook.ID,
			exampleUser.ID,
		).Return(errors.New("blah"))
		s.webhookDataManager = wd

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, wd)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
