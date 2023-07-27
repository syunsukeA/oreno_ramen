package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RedirectHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/signin")
}

func HomeHandlerGET(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func SigninHandlerGET(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

func SigninHandlerPOST(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)

	// ユーザ名とパスワードをDB検索

	c.Redirect(http.StatusFound, "/home")
}

func SignupHandlerGET(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func SignupHandlerPOST(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println("New user registration:")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)

	// RedirectHandler(c)
	c.Redirect(http.StatusFound, "/signin")
}

func SignoutHandlerGET(c *gin.Context) {
	// セッションのユーザ名を空に設定してセッション状態に保存するコードが必要
	c.Redirect(http.StatusFound, "/signin")
}
