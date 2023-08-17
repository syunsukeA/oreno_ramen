package dao

import (
	"log"
	_"reflect"
	"errors"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/jmoiron/sqlx"
)

// repsitory.Reviewのimplements
type Review struct {
	DB *sqlx.DB
}

func (r *Review)GetUnvisitedReviews() (ROs []*object.Review){
	entity := new(object.Review)
	ROs = append(ROs, entity)
	return ROs
}

func (r *Review)AddReviewAndShop(c *gin.Context, shopID string, userID int64, shopname string, req *object.CreateReviewRequest) (ro *object.Review, err error) {
	ro = new(object.Review)
	// トランザクション処理の開始
	tx, err := r.DB.BeginTxx(c, nil)
	if err != nil {
        return nil, err
    }
	defer func() {
        switch r := recover(); {
		case r != nil:
            tx.Rollback()
        case err != nil:
            tx.Rollback()
		}
    }()
	so := new(object.Shop)
	// shopに該当データがあるか確認
	q := `SELECT * from shops WHERE shop_id = ?`
	err = tx.QueryRowxContext(c, q, shopID).StructScan(so)
	if err != nil {
		// shopになかったらshopデータの追加
		if err == sql.ErrNoRows {
			q = `INSERT INTO shops (shop_id, user_id, shopname) VALUES (?, ?, ?)`
			res, err := tx.Exec(q, shopID, userID, shopname)
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
	q = `INSERT INTO reviews (user_id, shop_id, shopname, dishname, content, evaluate, review_img) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(q, userID, shopID, shopname, req.DishName, req.Content, req.Evaluate, req.ReviewImg)
	if err != nil {
		return nil, err
	}
	// ToDo: 作成したreviewをResBodyに入れるために検索？？
	// review_id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// ToDo: もしかしたらこの辺の処理いらないかも？
	nrows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if nrows <= 0 {
		return nil, sql.ErrNoRows
	}
	// トランザクション処理をコミット
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ro, nil
}

func (r *Review)FindReviewsByShopID(c *gin.Context, userID int64, shopID string) (ros []*object.Review, err error){
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

func (r *Review)UpdateReview(c *gin.Context, ro *object.Review) (roPost *object.Review, err error) {
	roPost = new(object.Review)
	q := `
		UPDATE reviews
		SET shopname = ?, dishname = ?, content = ?, evaluate = ?, bookmark = ?, review_img = ?
		WHERE user_id = ? AND review_id = ?`
	res, err := r.DB.ExecContext(c, q, ro.ShopName, ro.DishName, ro.Content, ro.Evaluate, ro.Bookmark, ro.ReviewImg, ro.UserID, ro.ReviewID)
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

func (r *Review)RemoveReviewAndShop(c *gin.Context, userID int64, shopID string, reviewID int64) (ro *object.Review, err error) {
	ro = new(object.Review)
	// トランザクション処理の開始
	tx, err := r.DB.BeginTxx(c, nil)
	if err != nil {
        return nil, err
    }
	// defer内で戻り値の変更 (エラーハンドリング) はできないっぽいのでlogに吐き出す
	defer func() {
        switch r := recover(); {
		case r != nil:
            tx.Rollback()
        case err != nil:
            tx.Rollback()
		case ro == nil:
			tx.Rollback()
			log.Println(sql.ErrNoRows)
		}
    }()
	// Exexでreviewを削除
	q := `DELETE FROM reviews WHERE user_id = ? AND shop_id = ? AND review_id = ?`
	res, err := tx.ExecContext(c, q, userID, shopID, reviewID)
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
	// ToDo: user_id, shop_idで検索してreview件数が0ならshopsから削除
	q = `SELECT * FROM reviews WHERE user_id = ? AND shop_id = ?`
	res, err = tx.ExecContext(c, q, userID, shopID)
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
		res, err := tx.ExecContext(c, q, userID, shopID)
		if err != nil {
			return nil, err
		}
		// ToDo: resのRowAffectedでエラーハンドリングすべきかも？		
		_, err = res.RowsAffected()
		if err != nil {
			return nil, err
		}
	}
	// トランザクション処理をコミット
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ro, nil
}