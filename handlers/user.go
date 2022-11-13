package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adrisongomez/project-go/models"
	"github.com/adrisongomez/project-go/repository"
	"github.com/adrisongomez/project-go/server"
	"github.com/adrisongomez/project-go/utils"
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

const (
	HASH_COST = 8
)

type SignUpAndLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpAndMeResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"id_token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpAndLoginRequest{}
		error := json.NewDecoder(r.Body).Decode(&request)
		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, error := utils.HashText(request.Password)
		if error != nil {
			http.Error(w, "Error can be processed", http.StatusInternalServerError)
		}

		id, error := ksuid.NewRandom()
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Email:    request.Email,
			Password: *hashedPassword,
			Id:       id.String(),
		}

		err := repository.InsertUser(r.Context(), &user)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					http.Error(w, "Email is being used!", http.StatusForbidden)
					return
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			}
		}

		json.NewEncoder(w).Encode(SignUpAndMeResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpAndLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if !utils.ValidateHash(user.Password, request.Password) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := utils.GenerateToken(user, s.Config().JwtSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})

	}
}

func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(utils.USER_KEY).(*models.User)
		json.NewEncoder(w).Encode(SignUpAndMeResponse{
			Id:    user.Id,
			Email: user.Email,
		})

	}
}
