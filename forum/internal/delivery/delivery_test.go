package delivery

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginWithCorrectCredentials(t *testing.T) {
	// создание нового http-запрос для входа в систему с правильными учетными данными
	req, err := http.NewRequest("POST", "/login", strings.NewReader("email=test@example.com&password=testpassword"))
	if err != nil {
		t.Fatal(err)
	}

	// создание нового http-рекордера для записи ответа
	rr := httptest.NewRecorder()

	// вызов функции handler входа в систему
	http.HandlerFunc(Login).ServeHTTP(rr, req)

	// проверка кода состояния ответа
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Login handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	// проверка установки cookie сессии
	cookie := rr.Result().Cookies()[0]
	if cookie.Name != COOKIE_NAME {
		t.Errorf("Login handler did not set the session cookie")
	}
}

func TestLoginWithIncorrectCredentials(t *testing.T) {
	// создание нового http-запроса для входа в систему с неверными учетными данными
	req, err := http.NewRequest("POST", "/login", strings.NewReader("email=test@example.com&password=wrongpassword"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(Login).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Login handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) > 0 && cookies[0].Name == COOKIE_NAME {
		t.Errorf("Login handler set the session cookie with incorrect credentials")
	}
}

func TestLogout(t *testing.T) {
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// установка сессионный файл cookie в запросе
	cookie := &http.Cookie{Name: COOKIE_NAME, Value: "testsessionid"}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()

	http.HandlerFunc(Logout).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Logout handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	// проверка очистки куки-файла сессии
	cookies := rr.Result().Cookies()
	if len(cookies) > 0 && cookies[0].Value != "" {
		t.Errorf("Logout handler did not clear the session cookie")
	}
}
func TestPostPageLoadsSuccessfully(t *testing.T) {
	req, err := http.NewRequest("GET", "/post", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Post)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUserRedirectedToReadPage(t *testing.T) {
	form := url.Values{}
	form.Add("like", "like")
	form.Add("id", "1")
	req, err := http.NewRequest("POST", "/like", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
	if location := rr.Header().Get("Location"); location != "/read/?id=1" {
		t.Errorf("handler returned wrong redirect location: got %v want %v",
			location, "/read/?id=1")
	}
}

func TestUserCanLikePost(t *testing.T) {
	form := url.Values{}
	form.Add("like", "like")
	form.Add("id", "1")
	req, err := http.NewRequest("POST", "/like", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	// Verify that the post was liked successfully
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	row := database.QueryRow("select * from likes where postID = 1 and owner = 'testuser'")
	var id int
	var postID int
	var owner string
	var likes int
	var dislikes int
	row.Scan(&id, &postID, &owner, &likes, &dislikes)
	if likes != 1 || dislikes != 0 {
		t.Errorf("post was not liked successfully")
	}
}

func TestUserCanDislikePost(t *testing.T) {
	form := url.Values{}
	form.Add("dislike", "dislike")
	form.Add("id", "1")
	req, err := http.NewRequest("POST", "/like", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	// Check if the post is disliked by the user
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	rows, _ := db.Query("SELECT dislike FROM likes WHERE postID = 1 AND owner = 'testuser'")
	defer rows.Close()

	var disliked int
	for rows.Next() {
		rows.Scan(&disliked)
		if disliked != 1 {
			t.Errorf("expected post to be disliked, but it wasn't")
		}
	}
}
