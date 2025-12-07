package clean_cart_handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	httpPkg "github.com/jva44ka/ozon-simulator-go-cart/pkg/http"
)

type CartService interface {
	RemoveAllProducts(ctx context.Context, userId uuid.UUID) error
}

type CleanCartHandler struct {
	cartService CartService
}

func NewCleanCartHandler(cartService CartService) *CleanCartHandler {
	return &CleanCartHandler{cartService: cartService}
}

// @Summary      Очистить корзину пользователя
// @Description  Метод полностью очищает корзину пользователя.
// Если у пользователя нет корзины или она пуста, то, как и при успешной очистке корзины, необходимо вернуть код ответа 204 No Content.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        user_id  path  string  true  "Токен пользователя"
// @Success      200  {object}  CleanCartResponse
// @Failure      404  {object}  httpPkg.ErrorResponse
// @Router       /user/{user_id}/cart [delete]
func (h *CleanCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userIdRaw := r.PathValue("user_id")
	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		if err = httpPkg.NewErrorResponse(w, http.StatusBadRequest, "user_id must be valid uuid"); err != nil {
			fmt.Println("json.Encode failed ", err)

			return
		}

		return
	}

	err = h.cartService.RemoveAllProducts(r.Context(), userId)
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
