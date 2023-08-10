package object

import (
	"time"
)

type Shop struct {
	ShopID int64 `json:"shop_id" db:"shop_id"`
	UserID int64 `json:"user_id" db:"user_id"`
	ShopName string `json:"shopname" db:"shopname"`
	ShopImg string `json:"shop_img"`
	LatestContent string `json:"latest_content"`
	AvgEvaluate uint `json:"avg_evaluate"`
	Bookmark bool `json:"bookmark" db:"bookmark"`
	Reviews []*Review `json:"reviews"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}