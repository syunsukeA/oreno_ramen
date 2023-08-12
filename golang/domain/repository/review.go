package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type Review interface {
	GetUnvisitedReviews() (ROs []*object.Review)
	AddReviewAndShop(c *gin.Context, shopID string, userID int64, shopname string, content string, eval uint, review_img string) (ro *object.Review, err error)
}