package get_cart_items_by_user_id_handler

import "github.com/google/uuid"

type GetReviewsResponse struct {
	CartItems []CartItemResponse `json:"cart_items"`
}

type CartItemResponse struct {
	Id     uint64    `json:"id"`
	SkuId  uint64    `json:"sku_id"`
	UserId uuid.UUID `json:"user_id"`
	Count  uint32    `json:"count"`
}
