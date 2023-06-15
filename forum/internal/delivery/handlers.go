package delivery

import (
	"log"
	"net/http"
)

func Handlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/register", Register)
	mux.HandleFunc("/signin", handleLogin)
	mux.HandleFunc("/callback", handleCallback)
	mux.HandleFunc("/login-call", handleCallLogin)
	mux.HandleFunc("/post/", Post)
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/terms", Terms)
	mux.HandleFunc("/read/", Read)
	mux.HandleFunc("/like/", Like)
	mux.HandleFunc("/likecomment/", LikeComment)
	mux.HandleFunc("/deletepost", DeletePost)

	fileServer := http.FileServer(http.Dir("./ui/style"))
	mux.Handle("/style", http.NotFoundHandler())
	mux.Handle("/style/", http.StripPrefix("/style", fileServer))
	log.Println("Server listening on http://127.0.0.1:8000")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
