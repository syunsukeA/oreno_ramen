package dao

import (
	"database/sql"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
)

// repsitory.Reviewのimplements
type Shop struct {
	DB *sql.DB
}

func (r *Shop)GetUnvisitedShops() (SOs []*object.Shop){
	entity := new(object.Shop)
	SOs = append(SOs, entity)
	return SOs
}