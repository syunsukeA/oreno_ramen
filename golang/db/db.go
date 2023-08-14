package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

func ConnectDB() {
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

	DB, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
