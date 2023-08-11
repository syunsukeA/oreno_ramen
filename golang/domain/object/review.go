package object

import (
	"time"
)

type Review struct {
	ReviewID int64 `json:"review_id" db:"review_id"`
	UserID int64 `json:"user_id" db:"user_id"`
	ShopID string `json:"shop_id" db:"shop_id"`
	ShopName string `json:"shopname" db:"shopname"`
	Content string `json:"content" db:"content"`
	Evaluate uint `json:"evaluate" db:"evaluate"`
	ReviewImg string `json:"review_img" db:"review_img"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}