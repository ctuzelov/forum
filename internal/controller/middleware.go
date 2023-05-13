package controller

import (
	"context"
	"errors"
	"net/http"

	"forum/internal/models"
)

type contextKey string

const ctxKey contextKey = "data"

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

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
