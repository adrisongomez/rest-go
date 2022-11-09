package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/adrisongomez/project-go/handlers"
	"github.com/adrisongomez/project-go/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading enviroment variables")
	}

	PORT := os.Getenv("PORT")
	SECRET := os.Getenv("JWT_SECRET")
	DB_URL := os.Getenv("DATABASE_URL")

	s, error := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		DatabaseURL: DB_URL,
		JwtSecret:   SECRET,
	})

	if error != nil {
		log.Fatal(error)
	}

	bindRoutes := func(s server.Server, r *mux.Router) {
		r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
		r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
        r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	}

	s.Start(bindRoutes)
}
