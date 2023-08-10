package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
)

type Shop interface {
	GetUnvisitedShops() (SOs []*object.Shop) // DB操作コマンド名考える。Selectとか...？
}