package delivery

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterGetMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", ct, "text/html; charset=utf-8")
	}

	expected := "<html><body><h1>Register</h1><form action=\"/register\" method=\"post\"><label for=\"email\">Email:</label><input type=\"email\" id=\"email\" name=\"email\"><br><input type=\"submit\" value=\"Submit\"></form></body></html>"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRegisterValidPostMethod(t *testing.T) {
	body := bytes.NewBufferString("email=testuser@example.com&username=testuser&password=testpassword")
	req, err := http.NewRequest("POST", "/register", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify user is inserted in database
	database, _ := sql.Open("sqlite3", "./forum.db")
	defer database.Close()
	rows, _ := database.Query("select * from users where email ='testuser@example.com'")
	defer rows.Close()
	var id int
	var username string
	var email string
	var password string
	for rows.Next() {
		rows.Scan(&id, &email, &username, &password)
	}
	if email != "testuser@example.com" || username != "testuser" || password == "testpassword" {
		t.Errorf("handler failed to insert new user")
	}
}

func TestRegisterInvalidPostMethod(t *testing.T) {
	body := bytes.NewBufferString("email=invalidemail&username=testuser&password=testpassword")
	req, err := http.NewRequest("POST", "/register", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "change your username"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
