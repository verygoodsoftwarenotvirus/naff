package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := httpRoutesTestDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock5 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItemList := fake.BuildFakeItemList()

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return((*v1.ItemList)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with error fetching items from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return((*v1.ItemList)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItemList := fake.BuildFakeItemList()

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})
}

func TestItemsService_SearchHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock3.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error conducting search", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)

		si := &mock3.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return([]uint64{}, errors.New("blah"))
		s.search = si

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si)
	})

	T.Run("with now rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock3.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock3.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, errors.New("blah"))
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock3.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})
}

func TestItemsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		mc := &mock4.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.itemCounter = mc

		r := &mock5.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r, si, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

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

	T.Run("with error creating item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		mc := &mock4.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.itemCounter = mc

		r := &mock5.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r, si, ed)
	})
}

func TestItemsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(true, nil)
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with no such item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(false, sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(false, errors.New("blah"))
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})
}

func TestItemsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with no such item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})
}

func TestItemsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock5.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, r, itemDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

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

	T.Run("with no rows fetching item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error updating item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock5.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, r, itemDataManager, ed)
	})
}

func TestItemsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock5.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Delete", mock1.Anything, exampleItem.ID).Return(nil)
		s.search = si

		mc := &mock4.UnitCounter{}
		mc.On("Decrement", mock1.Anything).Return()
		s.itemCounter = mc

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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r)
	})

	T.Run("with no item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_includeOwnerFetchers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := includeOwnerFetchers(proj, typ)

		expected := `
package main

import ()

func main() {

	s.userIDFetcher = userIDFetcher

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := includeOwnerFetchers(proj, typ)

		expected := `
package main

import ()

func main() {

	s.thingIDFetcher = thingIDFetcher
	s.anotherThingIDFetcher = anotherThingIDFetcher

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRelevantIDFetchers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRelevantIDFetchers(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		typ := proj.LastDataType()
		x := buildRelevantIDFetchers(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
)

func main() {
	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain that belongs to user", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range proj.DataTypes {
			proj.DataTypes[i].BelongsToAccount = true
			proj.DataTypes[i].RestrictedToAccountMembers = true
		}

		typ := proj.LastDataType()

		x := buildRelevantIDFetchers(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
)

func main() {
	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleThing := fake.BuildFakeThing()
	exampleThing.BelongsToAccount = exampleUser.ID
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	exampleAnotherThing.BelongsToAccount = exampleUser.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_setupDataManagersForCreation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		indexToReturnFalse := 0
		returnErr := true
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := setupDataManagersForCreation(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr)

		expected := `
package main

import (
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	itemDataManager := &mock.ItemDataManager{}
	itemDataManager.On().Return()
	s.itemDataManager = itemDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		indexToReturnFalse := 0
		returnErr := true

		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := setupDataManagersForCreation(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr)

		expected := `
package main

import (
	"errors"
	mock1 "github.com/stretchr/testify/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	thingDataManager := &mock.ThingDataManager{}
	thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, errors.New("blah"))
	s.thingDataManager = thingDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain returning false, with no error", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		indexToReturnFalse := 0
		returnErr := false

		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := setupDataManagersForCreation(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr)

		expected := `
package main

import (
	mock1 "github.com/stretchr/testify/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	thingDataManager := &mock.ThingDataManager{}
	thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(false, nil)
	s.thingDataManager = thingDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain and early return", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		returnErr := false

		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		indexToReturnFalse := len(proj.FindOwnerTypeChain(typ)) - 1

		x := setupDataManagersForCreation(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr)

		expected := `
package main

import (
	mock1 "github.com/stretchr/testify/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	thingDataManager := &mock.ThingDataManager{}
	thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
	s.thingDataManager = thingDataManager

	anotherThingDataManager := &mock.AnotherThingDataManager{}
	anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(false, nil)
	s.anotherThingDataManager = anotherThingDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_determineMockExpecters(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		indexToStopAt := 1
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := []string{"itemDataManager"}
		actual := determineMockExpecters(proj, typ, indexToStopAt)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceListFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceListFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"testing"
)

func TestItemsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItemList := fake.BuildFakeItemList()

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return((*v1.ItemList)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with error fetching items from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return((*v1.ItemList)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItemList := fake.BuildFakeItemList()

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItems", mock1.Anything, exampleUser.ID, mock1.AnythingOfType("*models.QueryFilter")).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.ItemList")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceSearchFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceSearchFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsService_SearchHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock2.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error conducting search", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)

		si := &mock.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return([]uint64{}, errors.New("blah"))
		s.search = si

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si)
	})

	T.Run("with now rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock2.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock2.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, errors.New("blah"))
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		si := &mock.IndexManager{}
		si.On("Search", mock1.Anything, exampleQuery, exampleUser.ID).Return(exampleItemIDs, nil)
		s.search = si

		itemDataManager := &mock2.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with search disabled", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.SearchEnabled = false
		x := buildTestServiceSearchFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsService_SearchHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error conducting search", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si)
	})

	T.Run("with now rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return([]v1.Item{}, errors.New("blah"))
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleQuery := "whatever"
		exampleLimit := uint8(123)
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleItemIDs []uint64
		for _, x := range exampleItemList {
			exampleItemIDs = append(exampleItemIDs, x.ID)
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItemsWithIDs", mock1.Anything, exampleUser.ID, exampleLimit, exampleItemIDs).Return(exampleItemList, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("[]models.Item")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://todo.verygoodsoftwarenotvirus.ru?q=%s&limit=%d", exampleQuery, exampleLimit),
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.SearchHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, si, itemDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceCreateFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceCreateFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock5 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.itemCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock4.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock5.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r, si, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

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

	T.Run("with error creating item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("CreateItem", mock1.Anything, mock1.AnythingOfType("*models.ItemCreationInput")).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.itemCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock4.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock5.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r, si, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildTestServiceCreateFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestYetAnotherThingsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("CreateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThingCreationInput")).Return(exampleYetAnotherThing, nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.yetAnotherThingCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

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

	T.Run("with error creating yet another thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("CreateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThingCreationInput")).Return((*v1.YetAnotherThing)(nil), errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("CreateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThingCreationInput")).Return(exampleYetAnotherThing, nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.yetAnotherThingCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager, mc, r, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with forums app", func(t *testing.T) {
		proj := testprojects.BuildForumsApp()
		typ := proj.DataTypes[3]
		x := buildTestServiceCreateFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	exampleForum := fake.BuildFakeForum()
	forumIDFetcher := func(_ *http.Request) uint64 {
		return exampleForum.ID
	}

	exampleSubforum := fake.BuildFakeSubforum()
	exampleSubforum.BelongsToForum = exampleForum.ID
	subforumIDFetcher := func(_ *http.Request) uint64 {
		return exampleSubforum.ID
	}

	exampleThread := fake.BuildFakeThread()
	exampleThread.BelongsToSubforum = exampleSubforum.ID
	exampleThread.BelongsToAccount = exampleUser.ID
	threadIDFetcher := func(_ *http.Request) uint64 {
		return exampleThread.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.forumIDFetcher = forumIDFetcher
		s.subforumIDFetcher = subforumIDFetcher
		s.threadIDFetcher = threadIDFetcher
		s.userIDFetcher = userIDFetcher

		examplePost := fake.BuildFakePost()
		examplePost.BelongsToThread = exampleThread.ID
		examplePost.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakePostCreationInputFromPost(examplePost)

		forumDataManager := &mock.ForumDataManager{}
		forumDataManager.On("ForumExists", mock1.Anything, exampleForum.ID).Return(true, nil)
		s.forumDataManager = forumDataManager

		subforumDataManager := &mock.SubforumDataManager{}
		subforumDataManager.On("SubforumExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID).Return(true, nil)
		s.subforumDataManager = subforumDataManager

		threadDataManager := &mock.ThreadDataManager{}
		threadDataManager.On("ThreadExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID, exampleThread.ID).Return(true, nil)
		s.threadDataManager = threadDataManager

		postDataManager := &mock.PostDataManager{}
		postDataManager.On("CreatePost", mock1.Anything, mock1.AnythingOfType("*models.PostCreationInput")).Return(examplePost, nil)
		s.postDataManager = postDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.postCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Post")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, postDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.forumIDFetcher = forumIDFetcher
		s.subforumIDFetcher = subforumIDFetcher
		s.threadIDFetcher = threadIDFetcher
		s.userIDFetcher = userIDFetcher

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

	T.Run("with error creating post", func(t *testing.T) {
		s := buildTestService()

		s.forumIDFetcher = forumIDFetcher
		s.subforumIDFetcher = subforumIDFetcher
		s.threadIDFetcher = threadIDFetcher
		s.userIDFetcher = userIDFetcher

		examplePost := fake.BuildFakePost()
		examplePost.BelongsToThread = exampleThread.ID
		examplePost.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakePostCreationInputFromPost(examplePost)

		forumDataManager := &mock.ForumDataManager{}
		forumDataManager.On("ForumExists", mock1.Anything, exampleForum.ID).Return(true, nil)
		s.forumDataManager = forumDataManager

		subforumDataManager := &mock.SubforumDataManager{}
		subforumDataManager.On("SubforumExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID).Return(true, nil)
		s.subforumDataManager = subforumDataManager

		threadDataManager := &mock.ThreadDataManager{}
		threadDataManager.On("ThreadExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID, exampleThread.ID).Return(true, nil)
		s.threadDataManager = threadDataManager

		postDataManager := &mock.PostDataManager{}
		postDataManager.On("CreatePost", mock1.Anything, mock1.AnythingOfType("*models.PostCreationInput")).Return((*v1.Post)(nil), errors.New("blah"))
		s.postDataManager = postDataManager

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

		mock1.AssertExpectationsForObjects(t, postDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.forumIDFetcher = forumIDFetcher
		s.subforumIDFetcher = subforumIDFetcher
		s.threadIDFetcher = threadIDFetcher
		s.userIDFetcher = userIDFetcher

		examplePost := fake.BuildFakePost()
		examplePost.BelongsToThread = exampleThread.ID
		examplePost.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakePostCreationInputFromPost(examplePost)

		forumDataManager := &mock.ForumDataManager{}
		forumDataManager.On("ForumExists", mock1.Anything, exampleForum.ID).Return(true, nil)
		s.forumDataManager = forumDataManager

		subforumDataManager := &mock.SubforumDataManager{}
		subforumDataManager.On("SubforumExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID).Return(true, nil)
		s.subforumDataManager = subforumDataManager

		threadDataManager := &mock.ThreadDataManager{}
		threadDataManager.On("ThreadExists", mock1.Anything, exampleForum.ID, exampleSubforum.ID, exampleThread.ID).Return(true, nil)
		s.threadDataManager = threadDataManager

		postDataManager := &mock.PostDataManager{}
		postDataManager.On("CreatePost", mock1.Anything, mock1.AnythingOfType("*models.PostCreationInput")).Return(examplePost, nil)
		s.postDataManager = postDataManager

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.postCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Post")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, postDataManager, mc, r, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceExistenceFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceExistenceFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"testing"
)

func TestItemsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(true, nil)
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with no such item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(false, sql.ErrNoRows)
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ItemExists", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(false, errors.New("blah"))
		s.itemDataManager = itemDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildTestServiceExistenceFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"testing"
)

func TestYetAnotherThingsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("YetAnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(true, nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with no such yet another thing in database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("YetAnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(false, sql.ErrNoRows)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error fetching yet another thing from database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("YetAnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(false, errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceReadFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceReadFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"testing"
)

func TestItemsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})

	T.Run("with no such item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		s.itemDataManager = itemDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, itemDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildTestServiceReadFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	"net/http"
	"testing"
)

func TestYetAnotherThingsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(exampleYetAnotherThing, nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager, ed)
	})

	T.Run("with no such yet another thing in database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return((*v1.YetAnotherThing)(nil), sql.ErrNoRows)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error fetching yet another thing from database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return((*v1.YetAnotherThing)(nil), errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(exampleYetAnotherThing, nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceUpdateFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceUpdateFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestItemsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, r, itemDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

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

	T.Run("with no rows fetching item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error fetching item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return((*v1.Item)(nil), errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error updating item", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		exampleInput := fake.BuildFakeItemUpdateInputFromItem(exampleItem)

		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("GetItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(exampleItem, nil)
		itemDataManager.On("UpdateItem", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Index", mock1.Anything, exampleItem.ID, exampleItem).Return(nil)
		s.search = si

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.Item")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, r, itemDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildTestServiceUpdateFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestYetAnotherThingsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingUpdateInputFromYetAnotherThing(exampleYetAnotherThing)

		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(exampleYetAnotherThing, nil)
		yetAnotherThingDataManager.On("UpdateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(nil)
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

		mock1.AssertExpectationsForObjects(t, r, yetAnotherThingDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

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

	T.Run("with no rows fetching yet another thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingUpdateInputFromYetAnotherThing(exampleYetAnotherThing)

		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return((*v1.YetAnotherThing)(nil), sql.ErrNoRows)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error fetching yet another thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingUpdateInputFromYetAnotherThing(exampleYetAnotherThing)

		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return((*v1.YetAnotherThing)(nil), errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error updating yet another thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingUpdateInputFromYetAnotherThing(exampleYetAnotherThing)

		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(exampleYetAnotherThing, nil)
		yetAnotherThingDataManager.On("UpdateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, yetAnotherThingDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingUpdateInputFromYetAnotherThing(exampleYetAnotherThing)

		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("GetYetAnotherThing", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(exampleYetAnotherThing, nil)
		yetAnotherThingDataManager.On("UpdateYetAnotherThing", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.YetAnotherThing")).Return(errors.New("blah"))
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

		mock1.AssertExpectationsForObjects(t, r, yetAnotherThingDataManager, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_setupDataManagersForDeletion(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		indexToReturnFalse := 0
		returnErr := true
		returnFalse := true
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := setupDataManagersForDeletion(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr, returnFalse)

		expected := `
package main

import (
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	itemDataManager := &mock.ItemDataManager{}
	itemDataManager.On().Return()
	s.itemDataManager = itemDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		actualCallArgs := []jen.Code{}
		returnValues := []jen.Code{}
		indexToReturnFalse := 0
		returnErr := true
		returnFalse := true

		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := setupDataManagersForDeletion(proj, typ, actualCallArgs, returnValues, indexToReturnFalse, returnErr, returnFalse)

		expected := `
package main

import (
	"errors"
	mock1 "github.com/stretchr/testify/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
)

func main() {
	thingDataManager := &mock.ThingDataManager{}
	thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(false, errors.New("blah"))
	s.thingDataManager = thingDataManager

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServiceArchiveFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestServiceArchiveFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestItemsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(nil)
		s.itemDataManager = itemDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		si := &mock3.IndexManager{}
		si.On("Delete", mock1.Anything, exampleItem.ID).Return(nil)
		s.search = si

		mc := &mock4.UnitCounter{}
		mc.On("Decrement", mock1.Anything).Return()
		s.itemCounter = mc

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

		mock1.AssertExpectationsForObjects(t, itemDataManager, mc, r)
	})

	T.Run("with no item in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(sql.ErrNoRows)
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToAccount = exampleUser.ID
		s.itemIDFetcher = func(req *http.Request) uint64 {
			return exampleItem.ID
		}

		itemDataManager := &mock.ItemDataManager{}
		itemDataManager.On("ArchiveItem", mock1.Anything, exampleItem.ID, exampleUser.ID).Return(errors.New("blah"))
		s.itemDataManager = itemDataManager

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

		mock1.AssertExpectationsForObjects(t, itemDataManager)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildTestServiceArchiveFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestYetAnotherThingsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleThing := fake.BuildFakeThing()
	thingIDFetcher := func(_ *http.Request) uint64 {
		return exampleThing.ID
	}

	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = exampleThing.ID
	anotherThingIDFetcher := func(_ *http.Request) uint64 {
		return exampleAnotherThing.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("ArchiveYetAnotherThing", mock1.Anything, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(nil)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

		r := &mock2.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mock3.UnitCounter{}
		mc.On("Decrement", mock1.Anything).Return()
		s.yetAnotherThingCounter = mc

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

		mock1.AssertExpectationsForObjects(t, thingDataManager, anotherThingDataManager, yetAnotherThingDataManager, mc, r)
	})

	T.Run("with nonexistent thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(false, nil)
		s.thingDataManager = thingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager)
	})

	T.Run("with error checking thing existence", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, errors.New("blah"))
		s.thingDataManager = thingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager)
	})

	T.Run("with nonexistent another thing", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(false, nil)
		s.anotherThingDataManager = anotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager, anotherThingDataManager)
	})

	T.Run("with error checking another thing existence", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, errors.New("blah"))
		s.anotherThingDataManager = anotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager, anotherThingDataManager)
	})

	T.Run("with no yet another thing in database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("ArchiveYetAnotherThing", mock1.Anything, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(sql.ErrNoRows)
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager, anotherThingDataManager, yetAnotherThingDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.thingIDFetcher = thingIDFetcher
		s.anotherThingIDFetcher = anotherThingIDFetcher

		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		s.yetAnotherThingIDFetcher = func(req *http.Request) uint64 {
			return exampleYetAnotherThing.ID
		}

		thingDataManager := &mock.ThingDataManager{}
		thingDataManager.On("ThingExists", mock1.Anything, exampleThing.ID).Return(true, nil)
		s.thingDataManager = thingDataManager

		anotherThingDataManager := &mock.AnotherThingDataManager{}
		anotherThingDataManager.On("AnotherThingExists", mock1.Anything, exampleThing.ID, exampleAnotherThing.ID).Return(true, nil)
		s.anotherThingDataManager = anotherThingDataManager

		yetAnotherThingDataManager := &mock.YetAnotherThingDataManager{}
		yetAnotherThingDataManager.On("ArchiveYetAnotherThing", mock1.Anything, exampleAnotherThing.ID, exampleYetAnotherThing.ID).Return(errors.New("blah"))
		s.yetAnotherThingDataManager = yetAnotherThingDataManager

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

		mock1.AssertExpectationsForObjects(t, thingDataManager, anotherThingDataManager, yetAnotherThingDataManager)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
