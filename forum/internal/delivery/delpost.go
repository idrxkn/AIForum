package delivery

import (
	"database/sql"
	"html/template"
	"net/http"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL parameter
	postID := r.URL.Query().Get("id")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is authorized to delete the post
	cook, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			Errors(w, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := cook.Value
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	rows, _ := database.Query("select * from sessions where session ='" + sessionToken + "'")
	var id int
	var user string
	var session string
	for rows.Next() {
		rows.Scan(&id, &user, &session)
	}
	if user == "" {
		tmp, err := template.ParseFiles("/Users/tikosch/Downloads/forum 3/ui/html/login.html")
		if err != nil {
			Errors(w, http.StatusInternalServerError)
			return
		}
		tmp.Execute(w, nil)
		return
	}

	// Delete the post from the database
	database, _ = sql.Open("sqlite3", "./forum.db")
	defer database.Close()

	DB, err := database.Prepare("DELETE FROM post WHERE id=$1 AND owner=$2")
	if err != nil {
		Errors(w, http.StatusInternalServerError)
		return
	}
	if _, err := DB.Exec(postID, user); err != nil {
		Errors(w, http.StatusInternalServerError)
		return
	}

	// Redirect the user to the home page
	http.Redirect(w, r, "/post", http.StatusFound)
}
