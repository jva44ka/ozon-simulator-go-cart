package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/model"
)

type CartRepository interface {
	AddCartItem(_ context.Context, cartItem model.CartItem) error
	GetCartItemsByUserId(_ context.Context, userId uuid.UUID) ([]model.CartItem, error)
}

type ProductService interface {
	GetProductBySku(ctx context.Context, sku uint64) (*model.Product, error)
}

type CartService struct {
	cartRepository CartRepository
	productService ProductService
}

func NewCartService(cartRepository CartRepository, productService ProductService) *CartService {
	return &CartService{cartRepository: cartRepository, productService: productService}
}

func (s *CartService) AddProduct(ctx context.Context, userId uuid.UUID, sku uint64, count uint32) error {
	if sku < 1 {
		return errors.New("sku must be greater than zero")
	}

	if userId == uuid.Nil {
		return errors.New("user_id must be not nil")
	}

	if count < 1 {
		return errors.New("count must be greater than zero")
	}

	_, err := s.productService.GetProductBySku(ctx, sku)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			return fmt.Errorf("productService.GetProductBySku: %w", err)
		}

		return err
	}

	cartItem := model.CartItem{
		UserId: userId,
		SkuId:  sku,
		Count:  count,
	}

	err = s.cartRepository.AddCartItem(ctx, cartItem)
	if err != nil {
		return fmt.Errorf("cartRepository.AddCartItem :%w", err)
	}

	return nil
}

func (s *CartService) GetItemsByUserId(ctx context.Context, userId uuid.UUID) ([]model.CartItem, error) {
	if userId == uuid.Nil {
		return nil, errors.New("userId must be not Nil")
	}

	reviews, err := s.cartRepository.GetCartItemsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("cartRepository.GetCartItemsByUserId :%w", err)
	}

	return reviews, nil
}
