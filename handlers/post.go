package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adrisongomez/project-go/models"
	"github.com/adrisongomez/project-go/repository"
	"github.com/adrisongomez/project-go/server"
	"github.com/adrisongomez/project-go/utils"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
}

type ListPostResponse struct {
	Page  uint64         `json:"page"`
	Posts []*models.Post `json:"posts"`
}

func GetPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(post)

	}
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(utils.CLAIMS_KEY).(*models.AppClaims)
		var postRequest = UpsertPostRequest{}

		if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := ksuid.NewRandom()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post := models.Post{
			Id:          id.String(),
			PostContent: postRequest.PostContent,
			UserId:      claims.UserId,
		}

		if err := repository.InsertPost(r.Context(), &post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		postMessage := models.WebsocketMessage{
			Type:    "Post_Created",
			Payload: post,
		}

		s.Hub().Broadcast(postMessage, nil)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&PostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})

	}

}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(utils.CLAIMS_KEY).(*models.AppClaims)
		params := mux.Vars(r)
		var postUpdate = UpsertPostRequest{}
		if err := json.NewDecoder(r.Body).Decode(&postUpdate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post := models.Post{
			Id:          params["id"],
			PostContent: postUpdate.PostContent,
			UserId:      claims.UserId,
		}

		if err := repository.UpdatePost(r.Context(), &post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&PostUpdateResponse{
			Message: fmt.Sprintf("`post_content` updated on post %s", post.Id),
		})
	}
}

func DeletePostHanlder(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(utils.CLAIMS_KEY).(*models.AppClaims)
		params := mux.Vars(r)

		if err := repository.DeletePost(r.Context(), params["id"], claims.UserId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&PostUpdateResponse{
			Message: fmt.Sprintf("Post %s has been deleted", params["id"]),
		})
	}
}

func ListPostHanlder(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		var page = uint64(0)
		var err error
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		posts, err := repository.ListPost(r.Context(), page)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&ListPostResponse{
			Posts: posts,
			Page:  page,
		})
	}
}
