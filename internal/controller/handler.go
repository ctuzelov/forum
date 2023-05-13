package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum/internal/repository"
	"forum/internal/service"
)

type Handler struct {
	ErrorLog  *log.Logger
	Service   *service.Service
	Tempcache *template.Template
}

func NewHandler(logger *log.Logger, DB *sql.DB) (*Handler, error) {
	repo := repository.New(DB)
	service := service.New(repo)
	tempcache, err := template.ParseGlob("ui/html/*.html")
	return &Handler{logger, service, tempcache}, err
}

func Routes(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ui/assets/", http.StripPrefix("/ui/assets/", http.FileServer(http.Dir("./ui/assets/"))))
	mux.HandleFunc("/", h.homepage)
	mux.HandleFunc("/signup", h.signup)
	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/posts/create", h.requireauth(h.postcreate))
	mux.HandleFunc("/logout", h.logout)
	mux.HandleFunc("/posts/", h.postview)
	mux.HandleFunc("/reaction", h.requireauth(h.reaction))
	mux.HandleFunc("/posts/created", h.requireauth(h.createdposts))
	mux.HandleFunc("/posts/liked", h.requireauth(h.likedposts))
	return h.middleware(mux)
}
