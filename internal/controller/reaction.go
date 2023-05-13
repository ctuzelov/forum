package controller

import (
	"net/http"
	"strconv"
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
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
	case "post":
		err = h.Service.Posts.React(id, userID, reaction)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
	default:
		h.errorpage(w, http.StatusBadRequest, nil)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
