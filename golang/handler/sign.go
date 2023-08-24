package handler

import (
	"fmt"
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
	r := c.Request
	w := c.Writer

	// リクエストボディからユーザー情報を取得
	user := new(object.User)
	user.UserName = r.FormValue("username")
	user.Password = r.FormValue("password")

	// ctxからimg_urlを取得
	filename, exists := c.Get("imgFilename")
	// imgURLがない場合はerr
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong on img-uploading")
		return
	}

	// 画像取得エンドポイントのURLを格納
	profileImg := fmt.Sprintf("img/%s", filename.(string))
	// TODO: エラーハンドリングの追加
	user.ProfileImg = profileImg

	// ユーザーを登録
	_, err := h.Ur.SignupByUsername(c, user.UserName, user.Password, user.ProfileImg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
