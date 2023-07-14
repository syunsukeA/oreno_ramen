package main

import (
  "fmt"
  "database/sql"
	"github.com/go-sql-driver/mysql"
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

func main() {
  db := connectDB()
  defer db.Close()

  query := "SELECT * FROM users"
  rows, err := db.Query(query)
  var user_id sql.NullInt64
  var name sql.NullString
  var password sql.NullString
  for rows.Next() {
    rows.Scan(&user_id, &name, &password)
    fmt.Println(user_id, name, password)
  }
  if err != nil {
    panic(err)
  }

	router := gin.Default()
	router.GET("/", internal.GetShoplist)
  router.Run("localhost:8080")
}