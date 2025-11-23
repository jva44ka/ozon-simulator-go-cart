package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/model"
)

type stubCartRepo struct {
	addFn func(ctx context.Context, item model.CartItem) error
	getFn func(ctx context.Context, userId uuid.UUID) ([]model.CartItem, error)
}

func (s *stubCartRepo) AddCartItem(ctx context.Context, item model.CartItem) error {
	return s.addFn(ctx, item)
}

func (s *stubCartRepo) GetCartItemsByUserId(ctx context.Context, userId uuid.UUID) ([]model.CartItem, error) {
	return s.getFn(ctx, userId)
}

type stubProductService struct {
	getFn func(ctx context.Context, sku uint64) (*model.Product, error)
}

func (s *stubProductService) GetProductBySku(ctx context.Context, sku uint64) (*model.Product, error) {
	return s.getFn(ctx, sku)
}
