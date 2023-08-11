package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// repsitory.Shopのimplements
type Shop struct {
	DB *sql.DB
}

func (r *Shop)GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string) (visitedIDs []string, err error){
	q := fmt.Sprintf(`SELECT shop_id FROM shops WHERE shop_id IN (?%s)`, strings.Repeat(",?", len(searchedIDs)-1))
	// QueryContextに可変引数として渡すためにany型に変換
	anytypeIDs := make([]interface{}, 0)
	for _, id := range searchedIDs {
		anytypeIDs = append(anytypeIDs, id)
	}
	var rows *sql.Rows
	rows, err = r.DB.QueryContext(ctx, q, anytypeIDs...)
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