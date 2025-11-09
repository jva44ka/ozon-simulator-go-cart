package model

import "github.com/google/uuid"

type CartItem struct {
	Id     uint64
	SkuId  uint64
	UserId uuid.UUID
	Count  uint32
}
