package dao

import (
	"database/sql"
	"errors"
	"fmt"
	_ "log"
	"strings"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// repsitory.Shopのimplements
type Shop struct {
	DB *sqlx.DB
}

func (r *Shop) GetVisitedShopIDs(ctx *gin.Context, searchedIDs []string, userID int64) (visitedIDs []string, err error) {
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

func (s *Shop) GetLatestShopByUserID(c *gin.Context, userID int64, num int64) (sos []*object.Shop, err error) {
	sos = []*object.Shop{} // お店のスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をnumで設定。
	q := `
	SELECT * FROM shops
	WHERE user_id = ?
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := s.DB.QueryxContext(c, q, userID, num)
	if err != nil {
		// エラーハンドリング
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // データがなければnilを返す
		}
		return nil, err
	}
	defer rows.Close()

	// 各行を読み込みながらスライスに追加
	for rows.Next() {
		so := new(object.Shop)
		if err := rows.StructScan(so); err != nil {
			return nil, err
		}
		sos = append(sos, so)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sos, nil
}

func (s *Shop) GetBookmarkShopByUserID(c *gin.Context, userID int64, num int64) (sos []*object.Shop, err error) {
	sos = []*object.Shop{} // お店のスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をnumで設定。
	q := `
	SELECT * FROM shops
	WHERE user_id = ? AND bookmark = 1
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := s.DB.QueryxContext(c, q, userID, num)
	if err != nil {
		// エラーハンドリング
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // データがなければnilを返す
		}
		return nil, err
	}
	defer rows.Close()

	// 各行を読み込みながらスライスに追加
	for rows.Next() {
		so := new(object.Shop)
		if err := rows.StructScan(so); err != nil {
			return nil, err
		}
		sos = append(sos, so)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sos, nil
}
