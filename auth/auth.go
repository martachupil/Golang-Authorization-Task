package auth

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"path"
	"strings"
	//"github.com/PuerkitoBio/goquery"
	"encoding/json"
)

type Instance struct {
	DB *sql.DB
}


func (i Instance) SignupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, path.Join("templates", "signup.html"))
		return
	}


	username := strings.TrimSpace(req.FormValue("username"))
	password := strings.TrimSpace(req.FormValue("password"))
	password_2 := req.FormValue("password_2")
	email := strings.TrimSpace(req.FormValue("email"))

	var errors [] string

	if username == "" {
		errors = append(errors, "empty username")
	}
	if password == "" {
		errors = append(errors, "empty password")
	}
	if email == "" {
		errors = append(errors, "empty email")
	}
	if (password_2 != password) {
		//log.Fatal("diff pass added")
	}

	if req.Method == "POST" && len(errors) > 0 {
		text, _ := json.Marshal(errors)
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(text))
		return
	}
	if len(errors) == 0 {
		var user string

		err := i.DB.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {

				http.Error(res, "Server error, unable to create your account." + err.Error(), 500)
				return
			}

			_, err = i.DB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
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
}

func (i Instance) LoginPage(res http.ResponseWriter, req *http.Request) {
if req.Method != "POST" {
http.ServeFile(res, req, path.Join("templates", "login.html"))
return
}

username := req.FormValue("username")
password := req.FormValue("password")

var databaseUsername string
var databasePassword string

err := i.DB.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

if err != nil {
http.Redirect(res, req, "/login", 301)
return
}

err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
if err != nil {
http.Redirect(res, req, "/login", 301)
return
}

res.Write([]byte("Hello " + databaseUsername))

}