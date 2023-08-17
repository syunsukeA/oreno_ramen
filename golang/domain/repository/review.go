package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

// 引数多すぎ問題のために構造体を定義
type AddReviewInput struct {
	ShopID   string
	UserID   int64
	ShopName string
	Request  *object.CreateReviewRequest
}

type Review interface {
	GetUnvisitedReviews() (ROs []*object.Review)
	AddReviewAndShop(c *gin.Context, rev_info *AddReviewInput) (ro *object.Review, err error)
	FindReviewsByShopID(c *gin.Context, userID int64, shopID string) (ros []*object.Review, err error)
	UpdateReview(c *gin.Context, roPre *object.Review) (roPost *object.Review, err error)
}
