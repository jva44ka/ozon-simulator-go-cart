package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/model"
	"sync"
	"sync/atomic"
)

type InMemoryCartItemRepository struct {
	storage map[uuid.UUID][]model.CartItem
	mx      sync.RWMutex

	idFactory atomic.Uint64
}

func NewCartItemRepository(cap int) *InMemoryCartItemRepository {
	return &InMemoryCartItemRepository{
		storage: make(map[uuid.UUID][]model.CartItem, cap),
	}
}

func (r *InMemoryCartItemRepository) AddProduct(ctx context.Context, userId uuid.UUID, sku uint64, count uint32) error {
	reviewID := r.idFactory.Add(1)
	review.ID = model.ReviewID(reviewID)

	r.mx.Lock()
	defer r.mx.Unlock()
	r.storage[review.Sku] = append(r.storage[review.Sku], review)

	return review, nil
}

func (r *InMemoryCartItemRepository) GetCartItemsByUserId(_ context.Context, sku model.Sku) ([]model.Review, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.storage[sku], nil
}
