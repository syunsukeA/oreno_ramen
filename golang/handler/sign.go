package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"
)

type HSign struct {
	Ur repository.User
}

func (h *HSign) SignupUser(c *gin.Context) {
	// r := c.Request
	w := c.Writer

	// リクエストボディからユーザー情報を取得
	var user object.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// ユーザーを登録
	_, err = h.Ur.SignupByUsername(c, user.UserName, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
