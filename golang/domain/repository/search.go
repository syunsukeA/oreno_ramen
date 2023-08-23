package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type Shop interface {
	GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string, userID int64) (visitedIDs []string, err error) // DB操作コマンド名考える。Selectとか...？
	GetLatestShopByUserID(c *gin.Context, userID int64, num int64) (sos []*object.Shop, err error)
}
