package delivery

import (
	"database/sql"
	"html/template"
	"net/http"
)

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL parameter
	postID := r.URL.Query().Get("id")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is authorized to edit the post
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

	// Get the post data from the request form
	title := r.FormValue("newtitle")
	content := r.FormValue("newcontent")
	typeTemp := r.Form["newcategory"]
	types := ""
	for _, categ := range typeTemp {
		types += categ + " "
	}
	image := r.FormValue("newimage")

	// Update the post in the database
	database, _ = sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	_, err = database.Exec("UPDATE post SET title=?, content=?, type=?, image=? WHERE id=? AND owner=?", title, content, types, image, postID, user)
	if err != nil {
		Errors(w, http.StatusInternalServerError)
		return
	}

	// Redirect the user to the post page
	http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
}
