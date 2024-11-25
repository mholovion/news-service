package routes

import (
	"net/http"
	"strings"

	"github.com/mholovion/news-service/controllers"
)

func RegisterRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", controllers.ReadPosts)
	http.HandleFunc("/new", controllers.CreatePost)

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			controllers.UpdatePost(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/delete") {
			controllers.DeletePost(w, r)
		} else {
			controllers.ReadPost(w, r)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	})
}
