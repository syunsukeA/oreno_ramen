package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// [!!注意!!] repositoryを追加したらmain.goに反映すること！！
type HSign struct {
	Ur repository.User
}

func (h *HSign) RegisterUser(c *gin.Context) {
	r := c.Request
	w := c.Writer

	// リクエストボディからユーザー情報を取得
	var user object.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// パスワードをハッシュ化して保存
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// ユーザーを登録
	_, err = h.Ur.RegisterByUsername(c, user.UserName) // ToDo: RegisterByUsername の実装
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HSign) LoginUser(c *gin.Context) {
	r := c.Request
	w := c.Writer

	// リクエストボディからログイン情報を取得
	var loginInfo object.User
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ユーザー名からユーザー情報を取得
	user, err := h.Ur.FindByUsername(c, loginInfo.UserName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// パスワードの照合
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HSign) LogoutUser(c *gin.Context) {
	// ToDo: ログアウト処理を実装

	c.Writer.WriteHeader(http.StatusNoContent)
}
