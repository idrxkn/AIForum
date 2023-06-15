package delivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssogolang *oauth2.Config

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ssogolang = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type Us struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json: "name"`
	Verif   bool   `json:"verified_email"`
	Picture string `json:"picture"`
}

var (
	us1         Us
	randomState = "random"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := ssogolang.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := ssogolang.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create token: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(string(content))

	if err := json.Unmarshal(content, &us1); err != nil {
		log.Fatal(err)
	}
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	fmt.Println(us1.Email)
	username := us1.Name
	email := us1.Email
	password := us1.Id
	crpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("error")
	}
	password = string(crpassword)
	DB, _ := database.Prepare(`Insert into users(email, username, password) values(?, ?, ?)`)

	DB.Exec(email, username, password)
	sessionId := InMemorySession.Init(username)
	cookie := &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessionId,
		Expires: time.Now().Add(15 * time.Minute),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login-call", http.StatusSeeOther)

	fmt.Fprintf(w, "Response: %s", content)
}

func handleCallLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Errors(w, http.StatusNotFound)
		return
	}
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()

	email := us1.Email
	rows, _ := database.Query("select * from users where email like '" + email + "'")
	var id int
	var username2 string
	var email2 string
	var password2 string
	for rows.Next() {
		rows.Scan(&id, &email2, &username2, &password2)
	}

	if username2 != "" && email2 != "" {

		sessionId := InMemorySession.Init(username2)
		Nick := InMemorySession.Get(sessionId)
		rows, _ := database.Query("select * from sessions where user ='" + Nick + "'")
		var id int
		var user string
		var session string
		for rows.Next() {
			rows.Scan(&id, &user, &session)
		}
		if Nick == user {
			DB, _ := database.Prepare("update sessions set session=? where user=?")
			DB.Exec(sessionId, Nick)
		} else {
			DB, _ := database.Prepare(`Insert into sessions(user,session) values(?,?)`)
			DB.Exec(Nick, sessionId)
		}
		cookie := &http.Cookie{
			Name:    COOKIE_NAME,
			Value:   sessionId,
			Expires: time.Now().Add(15 * time.Minute),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
