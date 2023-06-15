package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLikeDislikeWithNoCookie(t *testing.T) {
	req, err := http.NewRequest("POST", "/like?id=1&dislike=dislike", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLikeWithValidCookie(t *testing.T) {
	cookie := http.Cookie{Name: COOKIE_NAME, Value: "valid-session-token"}
	req, err := http.NewRequest("POST", "/like?id=1&like=like", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusFound)
	}
}

func TestLikeDislikeWithInvalidValues(t *testing.T) {
	cookie := http.Cookie{Name: COOKIE_NAME, Value: "valid-session-token"}
	req, err := http.NewRequest("POST", "/like?id=1&like=invalid&dislike=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Like)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}
}
