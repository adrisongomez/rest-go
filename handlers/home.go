package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adrisongomez/project-go/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func HomeHandler(server server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to my shitty golang program",
			Status:  "Connected",
		})
	}
}
