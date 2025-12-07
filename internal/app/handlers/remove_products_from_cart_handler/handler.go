package remove_products_from_cart_handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	httpPkg "github.com/jva44ka/ozon-simulator-go-cart/pkg/http"
)

type CartService interface {
	RemoveProduct(ctx context.Context, userId uuid.UUID, sku uint64) error
}

type RemoveProductsFromCartHandler struct {
	cartService CartService
}

func NewRemoveProductsFromCartHandler(cartService CartService) *RemoveProductsFromCartHandler {
	return &RemoveProductsFromCartHandler{cartService: cartService}
}

// @Summary      Удалить товар из корзины
// @Description  Метод полностью удаляет все количество товара из корзины пользователя.
// Если у пользователя вовсе нет данной позиции, то возвращается такой же ответ, как будто бы все позиции данного sku были успешно удалены
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        user_id  path  string  true  "Токен пользователя"
// @Param        sku_id   path  uint64  true  "SKU товара"
// @Success      200  {object}  RemoveProductsFromCartResponse
// @Failure      404  {object}  httpPkg.ErrorResponse
// @Router       /user/{user_id}/cart/{sku_id} [delete]
func (h *RemoveProductsFromCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	skuRaw := r.PathValue("sku_id")
	sku, err := strconv.Atoi(skuRaw)
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	if sku < 1 {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "sku must be more than zero"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	userIdRaw := r.PathValue("user_id")
	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "user_id must be valid uuid"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	err = h.cartService.RemoveProduct(r.Context(), userId, uint64(sku))
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusInternalServerError, err.Error()); err != nil {
			return
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}
