package controller

import (
	"errors"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) reaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorpage(w, http.StatusMethodNotAllowed, nil)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	if err := r.ParseForm(); err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	object := r.PostForm.Get("object")
	reaction := r.PostForm.Get("reaction")
	userID := data.User.ID
	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil || (reaction != "like" && reaction != "dislike") {
		h.errorpage(w, http.StatusBadRequest, nil)
		return
	}
	switch object {
	case "comment":
		err = h.Service.Comments.React(id, userID, reaction)
	case "post":
		err = h.Service.Posts.React(id, userID, reaction)

	default:
		h.errorpage(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		if errors.Is(err, models.ErrInvalidObjectId) {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
