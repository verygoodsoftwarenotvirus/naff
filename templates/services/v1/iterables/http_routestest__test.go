package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		t.Parallel()

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
}

func Test_buildRelevantIDFetchers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRelevantIDFetchers(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
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
}

func Test_setupDataManagersForCreation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
}

func Test_determineMockExpecters(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
}

func Test_buildTestServiceCreateFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
}

func Test_buildTestServiceExistenceFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
}

func Test_buildTestServiceReadFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
}

func Test_buildTestServiceUpdateFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
}

func Test_setupDataManagersForDeletion(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
}

func Test_buildTestServiceArchiveFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
		exampleItem.BelongsToUser = exampleUser.ID
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
}
