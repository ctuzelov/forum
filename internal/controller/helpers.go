package controller

import (
	"bytes"
	"net/http"

	"forum/internal/models"
)

type Data struct {
	User         models.User
	Content      any
	IsAuthorized bool
	IsEmpty      bool
	ErrMsgs      map[string]string
	Categories   []*string
}

type ErrorData struct {
	Status int
	Text   string
}

func (h *Handler) templaterender(w http.ResponseWriter, status int, page string, data any) {
	buf := new(bytes.Buffer)
	err := h.Tempcache.ExecuteTemplate(buf, page, data)
	if err != nil {
		h.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (h *Handler) errorpage(w http.ResponseWriter, status int, err error) {
	if err != nil {
		h.ErrorLog.Printf("server error: %v", err)
	}
	msg := http.StatusText(status)
	errdata := ErrorData{status, msg}
	h.templaterender(w, status, "errors.html", errdata)
}
