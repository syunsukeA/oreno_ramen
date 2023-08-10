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

type PhotoType struct {
	L string `json:"l"`
	M string `json:"m"`
	S string `json:"s"`
}

type Photo struct {
	Mobile PhotoType `json:"mobile"`
	PC PhotoType `json:"pc"`
}

type HPShop struct {
	Access string `json:"access"`
	Address string `json:"address"`
	Close string `json:"close"`
	ID string `json:"id"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Midinght string `json:"midinght"`
	MobileAccess string `json:"mobile_access"`
	Name string `json:"name"`
	Open string `json:"open"`
	Photo Photo `json:"photo"`
}

type HPResult struct {
	APIVersion string `json:"-"`
	ResultsAvailable int8 `json:"-"`
	ResultsReturuned string `json:"-"`
	ResultsStart int8 `json:"-"`
	Shop []*HPShop `json:"shop"`

}
type HPResponse struct {
	Result HPResult `json:"results"`
}