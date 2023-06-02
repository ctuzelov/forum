package controller

import (
	"context"
	"errors"
	"net/http"

	"forum/internal/models"
)

type contextKey string

const ctxKey contextKey = "data"

func (h *Handler) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		data := &Data{}
		switch err {
		case http.ErrNoCookie:
			data.User = models.User{}
		case nil:
			data.User, err = h.Service.Users.GetByToken(cookie.Value)
			if err != nil && !errors.Is(err, models.ErrNoRecord) {
				h.errorpage(w, http.StatusInternalServerError, err)
				return
			}
			if data.User.Token != nil {
				data.IsAuthorized = true
			}
		default:
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) requireauth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(ctxKey).(*Data)
		if !data.IsAuthorized {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
