package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/adrisongomez/project-go/repository"
	"github.com/adrisongomez/project-go/server"
	"github.com/adrisongomez/project-go/utils"
)

var (
	NO_AUTH_REQUIRED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_REQUIRED {
		if strings.Contains(route, p) {
			return false
		}
	}

	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			claims, err := utils.ValidateToken(tokenString, s.Config().JwtSecret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			user, err := repository.GetUserById(r.Context(), claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), utils.CLAIMS_KEY, claims)
			ctx = context.WithValue(ctx, utils.USER_KEY, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
