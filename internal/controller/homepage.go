package controller

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorpage(w, http.StatusNotFound, nil)
		return
	}
	data := r.Context().Value(ctxKey).(*Data)
	var err error
	data.Categories, err = h.Service.GetAllCategories()
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	query := r.URL.Query()
	var catID []int
	for _, value := range query["filter"] {
		number, err := strconv.Atoi(value)
		if err != nil {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		if number > 0 {
			catID = append(catID, number)
		}
	}
	if catID != nil {
		data.Content, err = h.Service.Posts.GetFiltered(catID)
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		data.Content, err = h.Service.Posts.GetAll()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
	}
	data.IsEmpty = (len(data.Content.([]*models.Post)) == 0)
	h.templaterender(w, http.StatusOK, "index.html", data)
}
