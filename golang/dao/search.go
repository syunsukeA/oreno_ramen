package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// repsitory.Shopのimplements
type Shop struct {
	DB *sqlx.DB
}

func (r *Shop)GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string, userID int64) (visitedIDs []string, err error){
	q := fmt.Sprintf(`SELECT shop_id, user_id FROM shops WHERE user_id = ? AND shop_id IN (?%s)`, strings.Repeat(",?", len(searchedIDs)-1))
	// QueryContextに可変引数として渡すためにany型に変換
	anytypeParams := make([]interface{}, 0)
	anytypeParams = append(anytypeParams, userID)
	for _, id := range searchedIDs {
		anytypeParams = append(anytypeParams, id)
	}
	var rows *sql.Rows
	rows, err = r.DB.QueryContext(ctx, q, anytypeParams...)
	if err != nil {
		return visitedIDs, err
		
	}
	// ToDo: 可変長引数とかで取って来れないかな？
	var id string
	for rows.Next() {
		rows.Scan(&id)
		visitedIDs = append(visitedIDs, id)
	}
	return visitedIDs, nil
}