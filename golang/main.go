package main

import (
  "fmt"
  "database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
  "github.com/gin-gonic/gin"
)

func connectDB(dbURL string) *sql.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	c := mysql.Config{
		DBName:    "oreno_ramen",
		User:      "root",
		Passwd:    "pokemon18782",
		Addr:      dbURL,
		// Net:       "tcp",
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
  dbURL := "Local"
  db := connectDB(dbURL)
  defer db.Close()

  query := "SELECT * FROM user"
  rows, err := db.Query(query)
  var user_id int
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