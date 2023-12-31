package main

import (
	"fmt"
	_ "net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/syunsukeA/oreno_ramen/golang/dao"
	"github.com/syunsukeA/oreno_ramen/golang/handler"
	"github.com/syunsukeA/oreno_ramen/golang/internal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func connectDB() *sqlx.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	c := mysql.Config{
		DBName:    "oreno_ramen_db",
		User:      "root",
		Passwd:    "passwd",
		Addr:      "db:3306",
		Net:       "tcp",
		ParseTime: true,
		// Collation: "utf8mb4_unicode_ci",
		Loc: jst,
	}
	db, err := sqlx.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}

const (
	port = 8080
)

func main() {
	db := connectDB()
	defer db.Close()

	// repositoryの作成
	sr := dao.Shop{DB: db}
	ur := dao.User{DB: db}
	rr := dao.Review{DB: db}

	hAuth := handler.HAuth{Ur: &ur}
	hImg := handler.HImg{Ur: &ur}
	// EndPointの定義 (ToDo: もう少し長くなりそうなら別関数に切り出してもいいかも？)
	rt := gin.Default()
	rt.Use(handler.DBTransactMiddleWare(db))
	rt.Use(cors.New(cors.Config{
		/*
		   ToDo
		    : CORSのAllowOrigins設定の見直し
		      : Dockerからホスト名引っ張ってくるしかないか...？
		*/
		// アクセス許可するオリジン
		AllowOrigins: []string{
			"*",
		},
	}))
	rt.GET("/", internal.GetShoplist) // ToDo: リダイレクト処理...？でも今回フロントから叩いてるからいらないか。
	// sign API
	hSign := handler.HSign{Ur: &ur}
	rt.POST("/signup", hImg.ProfileImgHandler(), hSign.SignupUser)
	//  profile API
	profRt := rt.Group("/profile")
	profRt.Use(hAuth.AuthenticationMiddleware())
	{
		h := handler.HUser{Rr: &rr, Ur: &ur}
		profRt.GET("", h.UserProfile)
		profRt.POST("/", internal.GetShoplist)
	}
	// search API
	searchRt := rt.Group("/search")
	searchRt.Use(hAuth.AuthenticationMiddleware())
	{
		h := handler.HSearch{Sr: &sr, Ur: &ur, Rr: &rr}
		searchRt.GET("/visited", h.SearchVisited)
		searchRt.GET("/unvisited", h.SearchUnvisited)
	}
	// review API
	reviewRt := rt.Group("/review")
	reviewRt.Use(hAuth.AuthenticationMiddleware())
	{
		h := handler.HReview{Rr: &rr, Ur: &ur}
		reviewRt.POST("", hImg.ReviewImgMiddleWare(), h.CreateReview)
		reviewRt.GET("/:review_id", internal.GetShoplist)
		reviewRt.POST("/:review_id", hImg.ReviewImgMiddleWare(), h.UpdateReview)
		reviewRt.DELETE("/:review_id", h.RemoveReview)
	}
	homeRt := rt.Group("/home")
	homeRt.Use(hAuth.AuthenticationMiddleware())
	{
		h := handler.HHome{Sr: &sr, Rr: &rr, Ur: &ur}
		homeRt.GET("/review", h.HomeReview)
		homeRt.GET("/review/bookmark", h.HomeBookmarkReview)
		homeRt.GET("/review/evaluate", h.HomeEvaluateReview)
		homeRt.GET("/review/period", h.HomePeriodReview)
		homeRt.GET("/shop", h.HomeShop)
		homeRt.GET("/shop/bookmark", h.HomeBookmarkShop)
		// homeRt.GET("/shop/evaluate", h.HomeEvaluateShop)
	}
	// image API
	imgRt := rt.Group("/img")
	// 認証あった方がいい気がしなくもないけど、URLの取得多分できないだろうしまぁいいかな。
	// imgRt.Use(hAuth.AuthenticationMiddleware())
	{
		imgRt.GET("/:filename", hImg.ShowImg)
	}
	rt.Run(fmt.Sprintf(":%d", port))
}
