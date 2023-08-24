package dao

import (
	"database/sql"
	"errors"
	"log"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type User struct {
	DB *sqlx.DB
}

func (r *User) FindByUsername(c *gin.Context, username string) (uo *object.User, err error) {
	uo = new(object.User)
	q := `SELECT * from users where username = ?`
	err = r.DB.QueryRowxContext(c, q, username).StructScan(uo)
	if err != nil {
		log.Println("=========================================")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return uo, nil
}

func (r *User) SignupByUsername(c *gin.Context, username string, password string, profileImg string) (uo *object.User, err error) {
	uo = new(object.User)

	// ユーザーをデータベースに登録
	q := `INSERT INTO users (username, password, profile_img) VALUES (?, ?, ?)`
	_, err = r.DB.ExecContext(c, q, username, password, profileImg)
	if err != nil {
		log.Println("=========================================")
		log.Println("Error inserting into users:", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// 登録したユーザーを取得して返す
	uo, err = r.FindByUsername(c, username)
	if err != nil {
		return nil, err
	}

	return uo, nil
}
