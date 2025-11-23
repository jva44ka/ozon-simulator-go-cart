package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/cart_items/service"
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

func TestCartService_AddProduct_OK(t *testing.T) {
	userId := uuid.New()

	cartRepo := &stubCartRepo{
		addFn: func(ctx context.Context, item model.CartItem) error {
			require.Equal(t, userId, item.UserId)
			require.Equal(t, uint64(10), item.SkuId)
			require.Equal(t, uint32(3), item.Count)
			return nil
		},
	}

	productSrv := &stubProductService{
		getFn: func(ctx context.Context, sku uint64) (*model.Product, error) {
			return &model.Product{Sku: sku}, nil
		},
	}

	svc := service.NewCartService(cartRepo, productSrv)

	err := svc.AddProduct(context.Background(), userId, 10, 3)
	require.NoError(t, err)
}

func TestCartService_AddProduct_InvalidSku(t *testing.T) {
	svc := service.NewCartService(nil, nil)

	err := svc.AddProduct(context.Background(), uuid.New(), 0, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "sku must be")
}

func TestCartService_AddProduct_InvalidUserId(t *testing.T) {
	svc := service.NewCartService(nil, nil)

	err := svc.AddProduct(context.Background(), uuid.Nil, 10, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "user_id")
}

func TestCartService_AddProduct_ProductNotFound(t *testing.T) {
	userId := uuid.New()

	productSrv := &stubProductService{
		getFn: func(ctx context.Context, sku uint64) (*model.Product, error) {
			return nil, model.ErrProductNotFound
		},
	}

	cartRepo := &stubCartRepo{}

	svc := service.NewCartService(cartRepo, productSrv)

	err := svc.AddProduct(context.Background(), userId, 10, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, model.ErrProductNotFound))
}

func TestCartService_AddProduct_RepoError(t *testing.T) {
	userId := uuid.New()

	productSrv := &stubProductService{
		getFn: func(ctx context.Context, sku uint64) (*model.Product, error) {
			return &model.Product{}, nil
		},
	}

	cartRepo := &stubCartRepo{
		addFn: func(ctx context.Context, item model.CartItem) error {
			return errors.New("db failure")
		},
	}

	svc := service.NewCartService(cartRepo, productSrv)

	err := svc.AddProduct(context.Background(), userId, 10, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cartRepository")
}

func TestCartService_GetItemsByUserId_OK(t *testing.T) {
	userId := uuid.New()

	expected := []model.CartItem{
		{UserId: userId, SkuId: 1, Count: 2},
	}

	cartRepo := &stubCartRepo{
		getFn: func(ctx context.Context, uid uuid.UUID) ([]model.CartItem, error) {
			require.Equal(t, userId, uid)
			return expected, nil
		},
	}

	svc := service.NewCartService(cartRepo, nil)

	items, err := svc.GetItemsByUserId(context.Background(), userId)
	require.NoError(t, err)
	require.Equal(t, expected, items)
}

func TestCartService_GetItemsByUserId_InvalidUser(t *testing.T) {
	svc := service.NewCartService(nil, nil)

	items, err := svc.GetItemsByUserId(context.Background(), uuid.Nil)
	require.Error(t, err)
	require.Nil(t, items)
}

func TestCartService_GetItemsByUserId_RepoError(t *testing.T) {
	userId := uuid.New()

	cartRepo := &stubCartRepo{
		getFn: func(ctx context.Context, uid uuid.UUID) ([]model.CartItem, error) {
			return nil, errors.New("storage error")
		},
	}

	svc := service.NewCartService(cartRepo, nil)

	items, err := svc.GetItemsByUserId(context.Background(), userId)
	require.Error(t, err)
	require.Nil(t, items)
	require.Contains(t, err.Error(), "cartRepository")
}
