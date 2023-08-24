package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DBTransactMiddleWare(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		// r := c.Request
		w := c.Writer
		tx, err := db.Beginx()
		if err != nil {
			log.Println("Transaction start error.")
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			c.Abort()
			return
		}
		c.Set("tx", tx)
		defer func() {
			// panicの検出
			if r := recover(); r != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
			}
			// StatusCodeの確認
			// ToDo: 正常コードの確認
			if c.Writer.Status() != http.StatusOK && c.Writer.Status() != http.StatusNoContent {
				// ロールバック
				tx.Rollback()
			} else {
				// コミット
				tx.Commit()
			}
		}()

		c.Next()
	}
}