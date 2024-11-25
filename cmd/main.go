package main

import (
	"log"
	"net/http"

	"github.com/mholovion/news-service/config"
	"github.com/mholovion/news-service/routes"
)

func main() {
	config.ConnectDB()

	routes.RegisterRoutes()

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
