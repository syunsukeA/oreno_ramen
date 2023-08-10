package repository

import (
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
)

type Review interface {
	GetUnvisitedReviews() (ROs []*object.Review)
}