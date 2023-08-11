package repository

import (
	_"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type Shop interface {
	GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string, userID int64) (visitedIDs []string, err error) // DB操作コマンド名考える。Selectとか...？
}