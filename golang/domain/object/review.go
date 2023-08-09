package object

import (
	"time"
)

type Review struct {
	ReviewID int64
	UserID int64
	ShopID int64
	ShopName string
	Content string
	Evaluate uint
	ReviewImg string
	Bookmark bool
	CreatedAt time.Time
	UpdatedAt time.Time
}