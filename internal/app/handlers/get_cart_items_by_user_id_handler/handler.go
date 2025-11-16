package get_cart_items_by_user_id_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jva44ka/ozon-simulator-go-cart/internal/domain/model"
	http2 "github.com/jva44ka/ozon-simulator-go-cart/pkg/http"
)

type CartService interface {
	GetItemsByUserId(ctx context.Context, userId uuid.UUID) ([]model.CartItem, error)
}

type GetReviewsBySkuHandler struct {
	cartService CartService
}

func NewGetCartItemsByUserIdHandler(cartService CartService) *GetReviewsBySkuHandler {
	return &GetReviewsBySkuHandler{cartService: cartService}
}

func (h *GetReviewsBySkuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIdRaw := r.PathValue("user_id")
	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		if err = http2.NewErrorResponse(w, http.StatusBadRequest, "user_id must be valid uuid"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	if userId == uuid.Nil {
		if err = http2.NewErrorResponse(w, http.StatusBadRequest, "userId must be not Nil"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	cartItems, err := h.cartService.GetItemsByUserId(r.Context(), userId)
	if err != nil {
		if err = http2.NewErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {
			return
		}

		return
	}

	response := GetReviewsResponse{CartItems: make([]CartItemResponse, 0, len(cartItems))}
	for _, cartItem := range cartItems {
		response.CartItems = append(response.CartItems, CartItemResponse{
			Id:     cartItem.Id,
			SkuId:  cartItem.SkuId,
			UserId: cartItem.UserId,
			Count:  cartItem.Count,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Println("success status failed")
		return
	}

	return
}
