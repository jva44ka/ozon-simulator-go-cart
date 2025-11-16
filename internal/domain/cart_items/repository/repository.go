package repository

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/model"
)

type InMemoryCartItemRepository struct {
	storage []model.CartItem
	mutex   sync.RWMutex

	idFactory atomic.Uint64
}

func NewCartItemRepository(cap int) *InMemoryCartItemRepository {
	return &InMemoryCartItemRepository{
		storage: make([]model.CartItem, cap),
	}
}

func (r *InMemoryCartItemRepository) AddCartItem(_ context.Context, cartItem model.CartItem) error {
	cartItemId := r.idFactory.Add(1)
	cartItem.Id = cartItemId

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, storageItem := range r.storage {
		if storageItem.SkuId == cartItem.SkuId && storageItem.UserId == cartItem.UserId {
			storageItem.Count = storageItem.Count + cartItem.Count

			return nil
		}
	}

	r.storage = append(r.storage, cartItem)

	return nil
}

func (r *InMemoryCartItemRepository) GetCartItemsByUserId(_ context.Context, userId uuid.UUID) ([]model.CartItem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]model.CartItem, 0, len(r.storage))

	for _, storageItem := range r.storage {
		if storageItem.UserId == userId {
			result = append(result, storageItem)
		}
	}

	return result, nil
}
