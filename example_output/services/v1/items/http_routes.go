package items

import (
	"database/sql"
	"net/http"
	"strconv"

	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"go.opencensus.io/trace"
)

var URIParamKey = "itemID"

func attachItemIDToSpan(span *trace.Span, itemID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("item_id", strconv.FormatUint(itemID, 10)))
	}
}

func attachUserIDToSpan(span *trace.Span, userID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("user_id", strconv.FormatUint(userID, 10)))
	}
}

// ListHandler is our list route
func (s *Service) ListHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ListHandler")
		defer span.End()
		qf := models.ExtractQueryFilter(req)
		userID := s.userIDFetcher(req)
		logger := s.logger.WithValue("user_id", userID)
		attachUserIDToSpan(span, userID)
		items, err := s.itemDatabase.GetItems(ctx, qf, userID)
		if err == sql.ErrNoRows {
			items = &models.ItemList{
				Items: []models.Item{},
			}
		} else if err != nil {
			logger.Error(err, "error encountered fetching items")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = s.encoderDecoder.EncodeResponse(res, items); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// CreateHandler is our item creation route
func (s *Service) CreateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "CreateHandler")
		defer span.End()
		userID := s.userIDFetcher(req)
		attachUserIDToSpan(span, userID)
		logger := s.logger.WithValue("user_id", userID)
		input, ok := ctx.Value(CreateMiddlewareCtxKey).(*models.ItemCreationInput)
		logger = logger.WithValue("input", input)
		if !ok {
			logger.Info("valid input not attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		input.BelongsTo = userID
		x, err := s.itemDatabase.CreateItem(ctx, input)
		if err != nil {
			s.logger.Error(err, "error creating item")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.itemCounter.Increment(ctx)
		attachItemIDToSpan(span, x.ID)
		s.reporter.Report(newsman.Event{
			Data: x,
			Topics: []string{
				topicName,
			},
			EventType: string(models.Create),
		})
		res.WriteHeader(http.StatusCreated)
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// ReadHandler returns a GET handler that returns an item
func (s *Service) ReadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ReadHandler")
		defer span.End()
		userID := s.userIDFetcher(req)
		itemID := s.itemIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id": userID,
			"item_id": itemID,
		})
		attachItemIDToSpan(span, itemID)
		attachUserIDToSpan(span, userID)
		x, err := s.itemDatabase.GetItem(ctx, itemID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error fetching item from database")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// UpdateHandler returns a handler that updates an item
func (s *Service) UpdateHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "UpdateHandler")
		defer span.End()
		input, ok := ctx.Value(UpdateMiddlewareCtxKey).(*models.ItemUpdateInput)
		if !ok {
			s.logger.Info("no input attached to request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		userID := s.userIDFetcher(req)
		itemID := s.itemIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"user_id": userID,
			"item_id": itemID,
		})
		attachItemIDToSpan(span, itemID)
		attachUserIDToSpan(span, userID)
		x, err := s.itemDatabase.GetItem(ctx, itemID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered getting item")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		x.Update(input)
		if err = s.itemDatabase.UpdateItem(ctx, x); err != nil {
			logger.Error(err, "error encountered updating item")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.reporter.Report(newsman.Event{
			Data: x,
			Topics: []string{
				topicName,
			},
			EventType: string(models.Update),
		})
		if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
			s.logger.Error(err, "encoding response")
		}
	}
}

// ArchiveHandler returns a handler that archives an item
func (s *Service) ArchiveHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "ArchiveHandler")
		defer span.End()
		userID := s.userIDFetcher(req)
		itemID := s.itemIDFetcher(req)
		logger := s.logger.WithValues(map[string]interface{}{
			"item_id": itemID,
			"user_id": userID,
		})
		attachItemIDToSpan(span, itemID)
		attachUserIDToSpan(span, userID)
		err := s.itemDatabase.ArchiveItem(ctx, itemID, userID)
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(err, "error encountered deleting item")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.itemCounter.Decrement(ctx)
		s.reporter.Report(newsman.Event{
			EventType: string(models.Archive),
			Data: &models.Item{
				ID: itemID,
			},
			Topics: []string{
				topicName,
			},
		})
		res.WriteHeader(http.StatusNoContent)
	}
}
