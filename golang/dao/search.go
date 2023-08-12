package dao

import (
	_"database/sql"
	"fmt"
	_"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// repsitory.Shopのimplements
type Shop struct {
	DB *sqlx.DB
}

func (r *Shop)GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string, userID int64) (visitedIDs []string, err error){
	q := fmt.Sprintf(`SELECT shop_id FROM shops WHERE user_id = ? AND shop_id IN (?%s)`, strings.Repeat(",?", len(searchedIDs)-1))
	// QueryContextに可変引数として渡すためにany型に変換
	anytypeParams := make([]interface{}, 0)
	anytypeParams = append(anytypeParams, userID)
	for _, id := range searchedIDs {
		anytypeParams = append(anytypeParams, id)
	}
	rows, err := r.DB.QueryxContext(ctx, q, anytypeParams...)
	if err != nil {
		return nil, err
	}
	// ToDo: 可変長引数とかで取って来れないかな？
	for rows.Next() {
		var str_id string
		if err := rows.Scan(&str_id); err != nil {
			return nil, err
		}
		visitedIDs = append(visitedIDs, str_id)
	}
	return visitedIDs, nil
}