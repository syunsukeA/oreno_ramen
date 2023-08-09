package main

import (
  "fmt"
  "database/sql"
	"github.com/go-sql-driver/mysql"
  _"net/http"
	"time"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
  "github.com/syunsukeA/oreno_ramen/golang/domain/object"
  "github.com/syunsukeA/oreno_ramen/golang/handler"

  "github.com/gin-gonic/gin"
)

func connectDB() *sql.DB {
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
		Loc:       jst,
	}
	db, err := sql.Open("mysql", c.FormatDSN())
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

  // objectの作成
  ro := object.Review{}

  // EndPointの定義 (ToDo: もう少し長くなりそうなら別関数に切り出してもいいかも？)
	rt := gin.Default()
  rt.GET("/", internal.GetShoplist)
  rt.POST("/signin", internal.GetShoplist)
  rt.POST("/signup", internal.GetShoplist)
  rt.POST("/signout", internal.GetShoplist)
  userRt := rt.Group("/user")
  {
    userRt.GET("/profile", internal.GetShoplist)
    userRt.POST("/profile", internal.GetShoplist)
  }
  searchRt := rt.Group("/search")
  {
    h := handler.HSearch{RO: ro}
    searchRt.GET("visited", internal.GetShoplist)
    searchRt.GET("unvisited", h.SearchUnvisited)
  }
  reviewRt := rt.Group("/review")
  {
    reviewRt.GET("reviews", internal.GetShoplist)
    reviewRt.POST("review", internal.GetShoplist)
    reviewRt.GET("review/{id}", internal.GetShoplist)
    reviewRt.POST("review/{id}", internal.GetShoplist)
    reviewRt.DELETE("review/{id}", internal.GetShoplist)
  }
  rt.Run(fmt.Sprintf(":%d", port))
}