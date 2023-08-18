package handler

import (
	"log"
	"net/http"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HAuth struct {
	Ur repository.User
}

/*
ToDo: 登録されていない場合、passとの照合がうまく行かない場合などのErrorHndlingについて考える。
あまり外部に知らせすぎるのもセキュリティ的に良くない...？
この辺は時間がある時に調べつつ考える。
*/
func (h *HAuth)AuthenticationMiddleware() gin.HandlerFunc {
	 return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		// ルーティング内の処理の前に実行される
		// Authorization Headerから読み出してbase64復号
		authUsername, authPassword, hasAuth := r.BasicAuth()
		// Authorization Headerがない場合、APIまで処理を行わずに401(Unauthorized)を返す
		if !hasAuth {
			log.Printf("Not-Authn")
			w.WriteHeader(http.StatusUnauthorized)
			c.Abort()
            return
		}
		log.Printf("username: %s, pass: %s", authUsername, authPassword)
		// DBにアクセスして登録済みuserか確認し、認証OKであればcontextにuser情報を格納
		var err error
		uo := new(object.User)
		uo, err = h.Ur.FindByUsername(c, authUsername)
		if err != nil {
			log.Printf("Internal server err")
			w.WriteHeader(http.StatusInternalServerError)
			c.Abort()
            return
		}
		// user名で検索して見つからない場合は処理の中断
		// ToDo: sign upにリダイレクトとかするべきなのかな...？
		if uo == nil {
			log.Printf("User not found")
			w.WriteHeader(http.StatusUnauthorized)
			c.Abort()
            return
		}
		// 登録passとHeaderのpassが一致しなければ Status NotAuth を返す
		// ToDo: passはハッシュ化してDBに保存した方がいいかも...？
		if uo.Password != authPassword {
			log.Printf("password is wrong")
			w.WriteHeader(http.StatusUnauthorized)
			c.Abort()
            return
		}
		// contextにuser情報を格納
		c.Set("authedUo", uo)
		c.Next() 
		// ルーティング内の処理の後に実行される
	} 
}