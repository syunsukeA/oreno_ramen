package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type Review interface {
	GetUnvisitedReviews() (ROs []*object.Review)
	AddReviewAndShop(c *gin.Context, shopID string, userID int64, shopname string, req *object.CreateReviewRequest) (ro *object.Review, err error)
	FindReviewsByShopID(c *gin.Context, userID int64, shopID string) (ros []*object.Review, err error)
	UpdateReview(c *gin.Context, roPre *object.Review) (roPost *object.Review, err error)
	RemoveReviewAndShop(c *gin.Context, userID int64, shopID string, reviewID int64) (ro *object.Review, err error)
}