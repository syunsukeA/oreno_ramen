package dao

import (
	"database/sql"
	"errors"
	"log"
	_ "reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/jmoiron/sqlx"
)

// repsitory.Reviewのimplements
type Review struct {
	DB *sqlx.DB
}

func (r *Review) FindByReviewID(c *gin.Context, reviewID int64) (ro *object.Review, err error) {
	ro = new(object.Review)
	q := `SELECT * from reviews where review_id = ?`
	err = r.DB.QueryRowxContext(c, q, reviewID).StructScan(ro)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return ro, nil
}

func (r *Review) GetLatestReviewByUserID(c *gin.Context, userID int64, num int64) (ros []*object.Review, err error) {
	ros = []*object.Review{} // レビューのスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をnumで設定。
	q := `
	SELECT * FROM reviews
	WHERE user_id = ?
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := r.DB.QueryxContext(c, q, userID, num)
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
		ro := new(object.Review)
		if err := rows.StructScan(ro); err != nil {
			return nil, err
		}
		ros = append(ros, ro)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ros, nil
}

func (r *Review) GetBookmarkReviewByUserID(c *gin.Context, userID int64, num int64) (ros []*object.Review, err error) {
	ros = []*object.Review{} // レビューのスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をnumで設定。
	q := `
	SELECT * FROM reviews
	WHERE user_id = ? AND bookmark = 1
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := r.DB.QueryxContext(c, q, userID, num)
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
		ro := new(object.Review)
		if err := rows.StructScan(ro); err != nil {
			return nil, err
		}
		ros = append(ros, ro)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ros, nil
}

func (r *Review) GetEvaluateReviewByUserID(c *gin.Context, userID int64, upper int64, lower int64) (ros []*object.Review, err error) {
	ros = []*object.Review{} // レビューのスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をupper、下限をlowerで設定。
	q := `
	SELECT * FROM reviews
	WHERE user_id = ? AND evaluate >= ? AND evaluate <= ?
	ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryxContext(c, q, userID, lower, upper)
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
		ro := new(object.Review)
		if err := rows.StructScan(ro); err != nil {
			return nil, err
		}
		ros = append(ros, ro)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ros, nil
}

func (r *Review) GetPeriodReviewByUserID(c *gin.Context, userID int64, upper time.Time, lower time.Time) (ros []*object.Review, err error) {
	ros = []*object.Review{} // レビューのスライスを初期化

	// SQLクエリの作成。userIDで絞り込み、作成日で降順にソートし、上限をupper、下限をlowerで設定。
	q := `
	SELECT * FROM reviews
	WHERE user_id = ? AND updated_at BETWEEN ? AND ?
	ORDER BY updated_at DESC
	`

	rows, err := r.DB.QueryxContext(c, q, userID, lower, upper)
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
		ro := new(object.Review)
		if err := rows.StructScan(ro); err != nil {
			return nil, err
		}
		ros = append(ros, ro)
	}

	// rows.Err()は、rowsの反復中にエラーが発生した場合にエラーを返す
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ros, nil
}

func (r *Review) GetUnvisitedReviews() (ROs []*object.Review) {
	entity := new(object.Review)
	ROs = append(ROs, entity)
	return ROs
}

func (r *Review) AddReviewAndShop(c *gin.Context, shopID string, userID int64, shopname string, req *object.CreateReviewRequest) (ro *object.Review, err error) {
	ro = new(object.Review)
	// トランザクション処理の開始
	any_tx, exists := c.Get("tx")
	// txがない場合は削除すべきものがないのでそのままreturn
	if !exists {
		return nil, sql.ErrConnDone // ToDo: このerrorは適当につけているなので後で適正なものを探そう...。
	}
	tx := any_tx.(*sqlx.Tx)
	so := new(object.Shop)
	// shopに該当データがあるか確認
	q := `SELECT * from shops WHERE shop_id = ?`
	err = tx.QueryRowxContext(c, q, shopID).StructScan(so)
	if err != nil {
		// shopになかったらshopデータの追加
		if err == sql.ErrNoRows {
			q = `INSERT INTO shops (shop_id, user_id, shopname, bookmark) VALUES (?, ?, ?, ?)`
			res, err := tx.Exec(q, shopID, userID, shopname, req.Bookmark)
			if err != nil {
				return nil, err
			}
			nrows, err := res.RowsAffected()
			if err != nil {
				return nil, err
			}
			if nrows <= 0 {
				return nil, sql.ErrNoRows
			}
		} else { // その他のerrorはシンプルに誤作動なのでerrをreturn
			return nil, err
		}
	}
	// reviewデータの追加
	// ToDo: 構造体駆使して短くする？
	q = `INSERT INTO reviews (user_id, shop_id, shopname, dishname, content, evaluate, bookmark, review_img) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(q, userID, shopID, shopname, req.DishName, req.Content, req.Evaluate, req.Bookmark, req.ReviewImg)
	if err != nil {
		return nil, err
	}
	// ToDo: 作成したreviewをResBodyに入れるために検索？？
	// review_id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	reviewID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	q = `SELECT * from reviews WHERE review_id = ?`
	err = tx.QueryRowxContext(c, q, reviewID).StructScan(ro)
	if err != nil {
		return nil, err
	}
	return ro, nil
}

func (r *Review) FindReviewsByShopID(c *gin.Context, userID int64, shopID string) (ros []*object.Review, err error) {
	ros = []*object.Review{}
	ro := new(object.Review)
	q := `SELECT * FROM reviews WHERE user_id = ? AND shop_id = ? LIMIT 20`
	rows, err := r.DB.QueryxContext(c, q, userID, shopID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(ro); err != nil {
			return nil, err
		}
		ros = append(ros, ro)
	}
	return ros, err
}

func (r *Review) UpdateReview(c *gin.Context, ro *object.Review) (roPost *object.Review, err error) {
	// トランザクション処理の開始
	any_tx, exists := c.Get("tx")
	// txがない場合は削除すべきものがないのでそのままreturn
	if !exists {
		return nil, sql.ErrConnDone // ToDo: このerrorは適当につけているなので後で適正なものを探そう...。
	}
	tx := any_tx.(*sqlx.Tx)
	roPost = new(object.Review)
	q := `
		UPDATE reviews
		SET shopname = ?, dishname = ?, content = ?, evaluate = ?, bookmark = ?, review_img = ?
		WHERE user_id = ? AND review_id = ?`
	res, err := tx.ExecContext(c, q, ro.ShopName, ro.DishName, ro.Content, ro.Evaluate, ro.Bookmark, ro.ReviewImg, ro.UserID, ro.ReviewID)
	if err != nil {
		return nil, err
	}
	n_affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	// 変更がなされなかった場合は両方nilを返す
	if n_affected == 0 {
		return nil, nil
	}
	// reviewを返すためにSELECTする
	q = `
		SELECT *
		FROM reviews
		WHERE user_id = ? AND review_id = ?`
	err = r.DB.QueryRowxContext(c, q, ro.UserID, ro.ReviewID).StructScan(roPost)
	if err != nil {
		return nil, err
	}
	log.Println("ro: ", ro)
	return roPost, nil
}

func (r *Review) RemoveReviewAndShop(c *gin.Context, ro *object.Review) (postRo *object.Review, err error) {
	// トランザクション処理の開始
	any_tx, exists := c.Get("tx")
	// txがない場合は削除すべきものがないのでそのままreturn
	if !exists {
		return nil, sql.ErrConnDone // ToDo: このerrorは適当につけているなので後で適正なものを探そう...。
	}
	tx := any_tx.(*sqlx.Tx)

	// Exexでreviewを削除
	q := `DELETE FROM reviews WHERE user_id = ? AND shop_id = ? AND review_id = ?`
	res, err := tx.ExecContext(c, q, ro.UserID, ro.ShopID, ro.ReviewID)
	if err != nil {
		return nil, err
	}
	// RowAffectedが0ならsql.ErrNoRowsを返す
	n_affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n_affected == 0 {
		return nil, nil
	}
	// user_id, shop_idで検索してreview件数が0ならshopsから削除
	q = `SELECT * FROM reviews WHERE user_id = ? AND shop_id = ?`
	res, err = tx.ExecContext(c, q, ro.UserID, ro.ShopID)
	if err != nil {
		return nil, err
	}
	n_affected, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n_affected == 0 {
		// shopsから削除
		q := `DELETE FROM shops WHERE user_id = ? AND shop_id = ?`
		res, err := tx.ExecContext(c, q, ro.UserID, ro.ShopID)
		if err != nil {
			return nil, err
		}
		// ToDo: resのRowAffectedでエラーハンドリングすべきかも？
		_, err = res.RowsAffected()
		if err != nil {
			return nil, err
		}
	}
	return ro, nil
}
