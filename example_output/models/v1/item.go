package models

import (
	"context"
	"net/http"
)

type (
	Item struct {
		ID         uint64
		Name       string
		Details    string
		CreatedOn  uint64
		UpdatedOn  *uint64
		ArchivedOn *uint64
		BelongsTo  uint64
	}
	ItemList struct {
		Pagination
		Items []Item
	}
	ItemCreationInput struct {
		Name      string
		Details   string
		BelongsTo uint64
	}
	ItemUpdateInput struct {
		Name      string
		Details   string
		BelongsTo uint64
	}
	ItemDataManager interface {
		GetItem(ctx context.Context, itemID, userID uint64) (*Item, error)
		GetItemCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllItemsCount(ctx context.Context) (uint64, error)
		GetItems(ctx context.Context, filter *QueryFilter, userID uint64) (*ItemList, error)
		GetAllItemsForUser(ctx context.Context, userID uint64) ([]Item, error)
		CreateItem(ctx context.Context, input *ItemCreationInput) (*Item, error)
		UpdateItem(ctx context.Context, updated *Item) error
		ArchiveItem(ctx context.Context, id, userID uint64) error
	}
	ItemDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler
		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an ItemInput with an Item
func (x *Item) Update(input *ItemUpdateInput) {
	if input.Name != "" || input.Name != x.Name {
		x.Name = input.Name
	}
	if input.Details != "" || input.Details != x.Details {
		x.Details = input.Details
	}
}
