package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/adrisongomez/project-go/databases"
	"github.com/adrisongomez/project-go/repository"
	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JwtSecret   string
	DatabaseURL string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JwtSecret == "" {
		return nil, errors.New("jwtSecret is required")
	}
	if config.DatabaseURL == "" {
		return nil, errors.New("DatabaseURL is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	repo, err := databases.NewPostgresRepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListAndSere: ", err)
	}
}
