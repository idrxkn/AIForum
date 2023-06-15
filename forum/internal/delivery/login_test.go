package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"forum/internal/delivery"
)

func TestLogin_ValidCredentials(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", strings.NewReader("email=test@example.com&password=password"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(delivery.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	cookie := rr.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, delivery.COOKIE_NAME) {
		t.Error("handler did not set session cookie")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", strings.NewReader("email=test@example.com&password=wrongpassword"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(delivery.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	cookie := rr.Header().Get("Set-Cookie")
	if cookie != "" {
		t.Error("handler set session cookie even with invalid credentials")
	}
}
