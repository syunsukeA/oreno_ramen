package dao

import (
	"database/sql"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
)

// repsitory.Reviewのimplements
type Review struct {
	DB *sql.DB
}

func (r *Review)GetUnvisitedReviews() (stROs []*object.Review){
	entity := new(object.Review)
	ROs = append(ROs, entity)
	return ROs
}