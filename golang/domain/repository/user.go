package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type User interface {
	FindByUsername(ctx *gin.Context, username string) (uo *object.User, err error)
}