package main

import (
	"database/sql"
	"fmt"
	// "github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	// "github.com/syunsukeA/oreno_ramen/golang/internal"
	"html/template"
	"net/http"
	"time"
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
		Loc: jst,
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

// func main() {
// 	db := connectDB()
// 	defer db.Close()

// 	query := "SELECT * FROM users"
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var user_id int
// 	var name string
// 	var password string
// 	var created_at time.Time
// 	for rows.Next() {
// 		rows.Scan(&user_id, &name, &password, &created_at)
// 		fmt.Println(user_id, name, password, created_at)
// 	}
// 	if err != nil {
// 		panic(err)
// 	}

// 	router := gin.Default()
// 	router.GET("/", internal.GetShoplist)
// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
// 		panic(err)
// 	}
// }

var store = sessions.NewCookieStore([]byte("secret-key"))

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
	"user3": "password3",
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/", redirectHandler)

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

	fmt.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println("Username:", username)
		fmt.Println("Password:", password)

		if pass, ok := users[username]; ok && pass == password {
			session, _ := store.Get(r, "session")
			session.Values["username"] = username
			session.Save(r, w)
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			renderTemplate(w, "login.html", "Invalid username or password")
		}
	} else {
		renderTemplate(w, "login.html", nil)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	renderTemplate(w, "index.html", username)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println("New user registration:")
		fmt.Println("Username:", username)
		fmt.Println("Password:", password)

		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		renderTemplate(w, "signup.html", nil)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = fmt.Sprintf("../templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
