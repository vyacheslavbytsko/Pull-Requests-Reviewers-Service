package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/api"
	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/db"
	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/handler"
)

func main() {
	database, dberr := db.NewDB()
	if dberr != nil {
		panic(dberr)
	}
	defer database.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	h := handler.NewHandler()
	apiHandler := api.Handler(h)

	router.Mount("/", apiHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	servererr := server.ListenAndServe()
	if servererr != nil {
		fmt.Println("Failed to start server:", servererr)
	}
}
