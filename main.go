package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"./auth"
)

var db *sql.DB
var err error

type Profile struct {
	Name    string
	Hobbies []string
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

	au := auth.Instance{DB:db}

	http.HandleFunc("/", index_handler)
	http.HandleFunc("/about/", about_handler)
	http.HandleFunc("/signup", au.SignupPage)
	http.HandleFunc("/login", au.LoginPage)

	fmt.Println(http.ListenAndServe(":8001", nil))
}
