package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

var db *sql.DB
var err error

type Profile struct {
	Name    string
	Hobbies []string
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, path.Join("templates", "signup.html"))
		return
	}

	username := strings.TrimSpace(req.FormValue("username"))
	password := strings.TrimSpace(req.FormValue("password"))
	password_2 := req.FormValue("password_2")
	email := strings.TrimSpace(req.FormValue("email"))

	if username == "" {
		log.Fatal("empty username")
	}
	if password == "" {
		log.Fatal("empty password")
	}
	if email == "" {
		log.Fatal("empty email")
	}
	if (password_2 != password) {
		log.Fatal("diff pass added")
	}

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {

			http.Error(res, "Server error, unable to create your account." + err.Error(), 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account." + err.Error(), 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account." + err.Error(), 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, path.Join("templates", "login.html"))
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}

func index_handler(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Guest", []string{"sport", "programming"}}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func about_handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/about/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "This site is about nothing")
}


func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}

func main() {

	// we create here an sql.DB and checking for errors
	db, err = sql.Open("mysql", "root:root@/mchsite")
	if err != nil {
		panic(err.Error())
	}
	// when func ends => close it
	defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}


	http.HandleFunc("/", index_handler)
	http.HandleFunc("/about/", about_handler)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.ListenAndServe(":8000", nil)
}