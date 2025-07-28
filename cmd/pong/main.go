package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aidenfine/pong/internal/api/status"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var PORT = ":8080"

func main() {
	_ = godotenv.Load()
	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		log.Fatal("NO DB URI")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Conected to DB")
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// timeout after 2 mins
	r.Use(middleware.Timeout(120 * time.Second))

	// mount routes
	r.Mount("/status", status.Routes(client))
	fmt.Println("Server started on port: ", PORT)

	log.Fatal(http.ListenAndServe(PORT, r))
}
