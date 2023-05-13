package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

func (h *Handler) postcreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := r.Context().Value(ctxKey).(*Data)
		h.templaterender(w, http.StatusOK, "create.html", data)
	case http.MethodPost:
		data := r.Context().Value(ctxKey).(*Data)
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("controller-postcreate-ParseForm: %w", err))
			return
		}

		var catid []int

		for _, value := range r.PostForm["cat"] {
			number, err := strconv.Atoi(value)
			if err != nil {
				h.errorpage(w, http.StatusBadRequest, err)
				return
			}
			catid = append(catid, number)
		}
		post := &models.Post{
			UserID:  data.User.ID,
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			CatID:   catid,
		}
		if post.Title == "" || post.Content == "" || post.CatID == nil {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		id, err := h.Service.Posts.Create(post)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
	}
}

func (h *Handler) postview(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil || len(path) != 3 || id < 1 {
		h.errorpage(w, http.StatusNotFound, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		post, err := h.Service.Posts.GetById(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				h.errorpage(w, http.StatusNotFound, err)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		post.Comments, err = h.Service.Comments.Fetch(post.ID)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		data.Content = post
		h.templaterender(w, http.StatusOK, "post.html", data)
	case http.MethodPost:
		if data.User == (models.User{}) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		comment := &models.Comment{
			PostID:  id,
			UserID:  data.User.ID,
			Content: r.PostForm.Get("content"),
		}
		comment.ParentID, err = strconv.Atoi(r.PostForm.Get("parent"))
		if comment.Content == "" || (err != nil && r.PostForm.Get("parent") != "") {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		err = h.Service.Comments.Create(comment)
		if err != nil {
			if errors.Is(err, models.ErrInvalidParent) {
				h.errorpage(w, http.StatusBadRequest, nil)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("postview-comments.create-%w", err))
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
	default:
		h.errorpage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) likedposts(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	posts, err := h.Service.Posts.GetUserLiked(data.User.ID)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("controller-likedposts-GetUserLiked: %w", err))
		return
	}
	data.IsEmpty = (len(posts) == 0)
	data.Content = posts
	h.templaterender(w, http.StatusOK, "likedcreated.html", data)
}

func (h *Handler) createdposts(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	posts, err := h.Service.Posts.GetUserCreated(data.User.ID)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, fmt.Errorf("controller-createdposts-GetUserLiked: %w", err))
		return
	}
	data.IsEmpty = (len(posts) == 0)
	data.Content = posts
	h.templaterender(w, http.StatusOK, "likedcreated.html", data)
}
