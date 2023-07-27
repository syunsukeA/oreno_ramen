package handler

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println("Username:", username)
		fmt.Println("Password:", password)

		// if pass, ok := users[username]; ok && pass == password {
		// 	session, _ := store.Get(r, "session")
		// 	session.Values["username"] = username
		// 	session.Save(r, w)
		// 	http.Redirect(w, r, "/home", http.StatusFound)
		// } else {
		// 	renderTemplate(w, "login.html", "Invalid username or password")
		// }
	}
}
