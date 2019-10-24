package items

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics"
)

const (
	CreateMiddlewareCtxKey models.ContextKey   = "item_create_input"
	UpdateMiddlewareCtxKey models.ContextKey   = "item_update_input"
	counterName            metrics.CounterName = "items"
	counterDescription                         = "the number of items managed by the items service"
	topicName              string              = "items"
	serviceName            string              = "items_service"
)

var _ models.ItemDataServer = (*Service)(nil)

type (
	Service struct {
		logger         logging.Logger
		itemCounter    metrics.UnitCounter
		itemDatabase   models.ItemDataManager
		userIDFetcher  UserIDFetcher
		itemIDFetcher  ItemIDFetcher
		encoderDecoder encoding.EncoderDecoder
		reporter       newsman.Reporter
	}
	UserIDFetcher func(*http.Request) uint64
	ItemIDFetcher func(*http.Request) uint64
)

// ProvideItemsService builds a new ItemsService
func ProvideItemsService(
	ctx context.Context,
	logger logging.Logger,
	db models.ItemDataManager,
	userIDFetcher UserIDFetcher,
	itemIDFetcher ItemIDFetcher,
	encoder encoding.EncoderDecoder,
	itemCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	itemCounter, err := itemCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}
	svc := &Service{
		logger:         logger.WithName(serviceName),
		itemDatabase:   db,
		encoderDecoder: encoder,
		itemCounter:    itemCounter,
		userIDFetcher:  userIDFetcher,
		itemIDFetcher:  itemIDFetcher,
		reporter:       reporter,
	}
	itemCount, err := svc.itemDatabase.GetAllItemsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current item count: %w", err)
	}
	svc.itemCounter.IncrementBy(ctx, itemCount)
	return svc, nil
}