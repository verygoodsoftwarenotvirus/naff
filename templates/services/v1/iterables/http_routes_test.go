package iterables

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := httpRoutesDotGo(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

const (
	// URIParamKey is a standard string that we'll use to refer to item IDs with.
	URIParamKey = "itemID"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := v1.ExtractQueryFilter(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// fetch items from database.
	items, err := s.itemDataManager.GetItems(ctx, userID, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		items = &v1.ItemList{
			Items: []v1.Item{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching items")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, items); err != nil {
		logger.Error(err, "encoding response")
	}
}

// SearchHandler is our search route.
func (s *Service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "SearchHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// we only parse the filter here because it will contain the limit
	filter := v1.ExtractQueryFilter(req)
	query := req.URL.Query().Get(v1.SearchQueryKey)
	logger = logger.WithValue("search_query", query)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	relevantIDs, searchErr := s.search.Search(ctx, query, userID)
	if searchErr != nil {
		logger.Error(searchErr, "error encountered executing search query")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fetch items from database.
	items, err := s.itemDataManager.GetItemsWithIDs(ctx, userID, filter.Limit, relevantIDs)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		items = []v1.Item{}
	} else if err != nil {
		logger.Error(err, "error encountered fetching items")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, items); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our item creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*v1.ItemCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)
	input.BelongsToUser = userID

	// create item in database.
	x, err := s.itemDataManager.CreateItem(ctx, input)
	if err != nil {
		logger.Error(err, "error creating item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachItemIDToSpan(span, x.ID)
	logger = logger.WithValue("item_id", x.ID)

	// notify relevant parties.
	s.itemCounter.Increment(ctx)
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(v1.Create),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "adding item to search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ExistenceHandler returns a HEAD handler that returns 200 if an item exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	tracing.AttachItemIDToSpan(span, itemID)
	logger = logger.WithValue("item_id", itemID)

	// fetch item from database.
	exists, err := s.itemDataManager.ItemExists(ctx, itemID, userID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking item existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// ReadHandler returns a GET handler that returns an item.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	tracing.AttachItemIDToSpan(span, itemID)
	logger = logger.WithValue("item_id", itemID)

	// fetch item from database.
	x, err := s.itemDataManager.GetItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching item from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates an item.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*v1.ItemUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)
	input.BelongsToUser = userID

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	logger = logger.WithValue("item_id", itemID)
	tracing.AttachItemIDToSpan(span, itemID)

	// fetch item from database.
	x, err := s.itemDataManager.GetItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update item in database.
	if err = s.itemDataManager.UpdateItem(ctx, x); err != nil {
		logger.Error(err, "error encountered updating item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(v1.Update),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "updating item in search index")
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ArchiveHandler returns a handler that archives an item.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	logger = logger.WithValue("item_id", itemID)
	tracing.AttachItemIDToSpan(span, itemID)

	// archive the item in the database.
	err = s.itemDataManager.ArchiveItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.itemCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(v1.Archive),
		Data:      &v1.Item{ID: itemID},
		Topics:    []string{topicName},
	})
	if indexDeleteErr := s.search.Delete(ctx, itemID); indexDeleteErr != nil {
		logger.Error(indexDeleteErr, "error removing item from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildHTTPRoutesConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildHTTPRoutesConstantDefs(typ)

		expected := `
package example

import ()

const (
	// URIParamKey is a standard string that we'll use to refer to item IDs with.
	URIParamKey = "itemID"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteLoggerAndTracingStatementsForListOfEntities(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteLoggerAndTracingStatementsForListOfEntities(proj, typ)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// determine user ID. userID := s . userIDFetcher (req)  tracing.AttachUserIDToSpan  (span,userID) logger = logger . WithValue ("user_id",userID)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteLoggerAndTracingStatementsForSingleEntity(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteLoggerAndTracingStatementsForSingleEntity(proj, typ)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// determine user ID. userID := s . userIDFetcher (req)  tracing.AttachUserIDToSpan  (span,userID) logger = logger . WithValue ("user_id",userID)
// determine item ID. itemID := s . itemIDFetcher (req)  tracing.AttachItemIDToSpan  (span,itemID) logger = logger . WithValue ("item_id",itemID)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSearchHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSearchHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// SearchHandler is our search route.
func (s *Service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "SearchHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// we only parse the filter here because it will contain the limit
	filter := v1.ExtractQueryFilter(req)
	query := req.URL.Query().Get(v1.SearchQueryKey)
	logger = logger.WithValue("search_query", query)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	relevantIDs, searchErr := s.search.Search(ctx, query, userID)
	if searchErr != nil {
		logger.Error(searchErr, "error encountered executing search query")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fetch items from database.
	items, err := s.itemDataManager.GetItemsWithIDs(ctx, userID, filter.Limit, relevantIDs)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		items = []v1.Item{}
	} else if err != nil {
		logger.Error(err, "error encountered fetching items")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, items); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildListHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildListHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// ensure query filter.
	filter := v1.ExtractQueryFilter(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// fetch items from database.
	items, err := s.itemDataManager.GetItems(ctx, userID, filter)
	if err == sql.ErrNoRows {
		// in the event no rows exist return an empty list.
		items = &v1.ItemList{
			Items: []v1.Item{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching items")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, items); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteLoggerAndTracingStatementsForModification(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		includeExistenceChecks := true
		includeSelf := true
		assignToUser := true
		assignToInput := true

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteLoggerAndTracingStatementsForModification(proj, typ, includeExistenceChecks, includeSelf, assignToUser, assignToInput)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// determine user ID. userID := s . userIDFetcher (req) logger = logger . WithValue ("user_id",userID)  tracing.AttachUserIDToSpan  (span,userID) input . BelongsToUser = userID
// determine item ID. itemID := s . itemIDFetcher (req) logger = logger . WithValue ("item_id",itemID)  tracing.AttachItemIDToSpan  (span,itemID)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// CreateHandler is our item creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check request context for parsed input struct.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*v1.ItemCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)
	input.BelongsToUser = userID

	// create item in database.
	x, err := s.itemDataManager.CreateItem(ctx, input)
	if err != nil {
		logger.Error(err, "error creating item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracing.AttachItemIDToSpan(span, x.ID)
	logger = logger.WithValue("item_id", x.ID)

	// notify relevant parties.
	s.itemCounter.Increment(ctx)
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(v1.Create),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "adding item to search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExistenceHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildExistenceHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// ExistenceHandler returns a HEAD handler that returns 200 if an item exists, 404 otherwise.
func (s *Service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ExistenceHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	tracing.AttachItemIDToSpan(span, itemID)
	logger = logger.WithValue("item_id", itemID)

	// fetch item from database.
	exists, err := s.itemDataManager.ItemExists(ctx, itemID, userID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err, "error checking item existence in database")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if exists {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildReadHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildReadHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// ReadHandler returns a GET handler that returns an item.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	tracing.AttachItemIDToSpan(span, itemID)
	logger = logger.WithValue("item_id", itemID)

	// fetch item from database.
	x, err := s.itemDataManager.GetItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching item from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// UpdateHandler returns a handler that updates an item.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// check for parsed input attached to request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*v1.ItemUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)
	input.BelongsToUser = userID

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	logger = logger.WithValue("item_id", itemID)
	tracing.AttachItemIDToSpan(span, itemID)

	// fetch item from database.
	x, err := s.itemDataManager.GetItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the data structure.
	x.Update(input)

	// update item in database.
	if err = s.itemDataManager.UpdateItem(ctx, x); err != nil {
		logger.Error(err, "error encountered updating item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.reporter.Report(newsman.Event{
		Data:      x,
		Topics:    []string{topicName},
		EventType: string(v1.Update),
	})
	if searchIndexErr := s.search.Index(ctx, x.ID, x); searchIndexErr != nil {
		logger.Error(searchIndexErr, "updating item in search index")
	}

	// encode our response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveHandlerFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveHandlerFuncDecl(proj, typ)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// ArchiveHandler returns a handler that archives an item.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine user ID.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// determine item ID.
	itemID := s.itemIDFetcher(req)
	logger = logger.WithValue("item_id", itemID)
	tracing.AttachItemIDToSpan(span, itemID)

	// archive the item in the database.
	err = s.itemDataManager.ArchiveItem(ctx, itemID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting item")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.itemCounter.Decrement(ctx)
	s.reporter.Report(newsman.Event{
		EventType: string(v1.Archive),
		Data:      &v1.Item{ID: itemID},
		Topics:    []string{topicName},
	})
	if indexDeleteErr := s.search.Delete(ctx, itemID); indexDeleteErr != nil {
		logger.Error(indexDeleteErr, "error removing item from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
