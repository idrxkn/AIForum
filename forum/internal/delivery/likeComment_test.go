package delivery

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type likeCommentTestCase struct {
	name           string
	cookie         *http.Cookie
	formValues     url.Values
	expectedStatus int
}

func TestLikeComment(t *testing.T) {
	cookie := &http.Cookie{Name: COOKIE_NAME, Value: "session_token"}

	testCases := []likeCommentTestCase{
		{
			name:           "Unauthorized user",
			cookie:         nil,
			formValues:     url.Values{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing 'id' parameter",
			cookie:         cookie,
			formValues:     url.Values{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Like a comment",
			cookie: cookie,
			formValues: url.Values{
				"likeC": {"like"},
			},
			expectedStatus: http.StatusFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/like-comment?id=1", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			req.Form = tc.formValues

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(LikeComment)
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
			}
		})
	}
}
