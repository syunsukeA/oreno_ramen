package main

import (
  "fmt"
  "database/sql"
	"github.com/go-sql-driver/mysql"
  "net/http"
	"time"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
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

  query := "SELECT * FROM users"
  rows, err := db.Query(query)
  if err != nil {
    panic(err)
  }
  var user_id int
  var name string
  var password string
  var created_at time.Time
  for rows.Next() {
    rows.Scan(&user_id, &name, &password, &created_at)
    fmt.Println(user_id, name, password, created_at)
  }
  if err != nil {
    panic(err)
  }

	router := gin.Default()
	router.GET("/", internal.GetShoplist)
  if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		panic(err)
	}
}