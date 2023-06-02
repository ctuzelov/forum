package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/models"
)

func (h *Handler) postcreate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	var err error
	data.Categories, err = h.Service.GetAllCategories()
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.templaterender(w, http.StatusOK, "create.html", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}

		var catid []int

		for _, value := range r.PostForm["cat"] {
			number, err := strconv.Atoi(value)
			if err != nil || number > 5 || number < 1 {
				h.errorpage(w, http.StatusBadRequest, nil)
				return
			}
			catid = append(catid, number)
		}
		post := &models.Post{
			UserID:    data.User.ID,
			Title:     strings.TrimSpace(r.PostForm.Get("title")),
			Content:   strings.TrimSpace(r.PostForm.Get("content")),
			CreatedAt: time.Now().Format("January 2, 2006"),
			CatID:     catid,
		}
		data.ErrMsgs = make(map[string]string)
		if len(post.Title) > 100 {
			post.Title = ""
			data.Content = post
			data.ErrMsgs["post"] = "Title is too long. Max - 100 symbols"
			h.templaterender(w, http.StatusBadRequest, "create.html", data)
			return
		}
		if post.CatID == nil {
			data.Content = post
			data.ErrMsgs["post"] = "Choose at least one category"
			h.templaterender(w, http.StatusBadRequest, "create.html", data)
			return
		}
		if post.Title == "" || post.Content == "" {
			data.Content = post
			data.ErrMsgs["post"] = "Please fill out all fields"
			h.templaterender(w, http.StatusBadRequest, "create.html", data)
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
		h.errorpage(w, http.StatusNotFound, nil)
		return
	}
	switch r.Method {
	case http.MethodGet:
		post, err := h.Service.Posts.GetById(id, data.User.ID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				h.errorpage(w, http.StatusNotFound, err)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		post.Comments, err = h.Service.Comments.Fetch(post.ID, data.User.ID)
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
			Content: strings.TrimSpace(r.PostForm.Get("content")),
		}
		comment.ParentID, err = strconv.Atoi(r.PostForm.Get("parent"))
		if comment.Content == "" || (err != nil && r.PostForm.Get("parent") != "") {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		_, err = h.Service.Posts.GetById(comment.PostID, comment.UserID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				h.errorpage(w, http.StatusBadRequest, nil)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		err = h.Service.Comments.Create(comment)
		if err != nil {
			if errors.Is(err, models.ErrInvalidParent) {
				h.errorpage(w, http.StatusBadRequest, nil)
				return
			}
			h.errorpage(w, http.StatusInternalServerError, err)
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
		h.errorpage(w, http.StatusInternalServerError, err)
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
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	data.IsEmpty = (len(posts) == 0)
	data.Content = posts
	h.templaterender(w, http.StatusOK, "likedcreated.html", data)
}
