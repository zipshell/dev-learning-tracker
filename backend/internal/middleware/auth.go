package mw

import (
	"context"
	"net/http"

	"github.com/zipshell/dev-learning-tracker/internal/auth"
)

func Auth(authSvc auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// validate auth (cookie, header, session, token, etc.)
			session, err := r.Cookie("session")
			if err != nil || session.Value == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			sessionUser, err := authSvc.ValidateSession(r.Context(), session.Value)
			if err != nil || session.Value == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", sessionUser)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
