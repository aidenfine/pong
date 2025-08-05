package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aidenfine/pong/internal/api/analytics"
	"github.com/aidenfine/pong/internal/api/status"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var PORT = ":8080"

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("duckdb", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Connected to DB")

	client := db

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	r.Mount("/status", status.Routes(client))
	r.Mount("/analytics", analytics.Routes(client))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server started on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
