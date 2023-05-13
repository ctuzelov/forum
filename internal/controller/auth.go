package controller

import (
	"net/http"

	"forum/internal/models"
	"forum/internal/validator"
)

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// data := r.Context().Value(ctxKey).(*Data)
		h.templaterender(w, http.StatusOK, "signup.html", nil)
	case http.MethodPost:
		data := r.Context().Value(ctxKey).(*Data)
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		form := models.User{
			Email:       r.PostForm.Get("email"),
			Name:        r.PostForm.Get("name"),
			Password:    r.PostForm.Get("password"),
			ConPassword: r.PostForm.Get("conpassword"),
		}
		data.Content = form
		if form.Email == "" || form.Name == "" || form.Password == "" || form.ConPassword == "" {
			data.ErrMsgs = validator.GetErrMsgs(form)
			h.templaterender(w, http.StatusBadRequest, "signup.html", data)
			return
		}
		data.ErrMsgs = validator.GetErrMsgs(form)
		if len(data.ErrMsgs) != 0 {
			h.templaterender(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}
		err = h.Service.SignUp(form)
		if err != nil {
			switch err {
			case models.ErrDuplicateEmail:
				data.ErrMsgs["email"] = validator.MsgEmailExists
			case models.ErrDuplicateName:
				data.ErrMsgs["name"] = validator.MsgNameExists
			default:
				h.errorpage(w, http.StatusInternalServerError, err)
				return
			}
			h.templaterender(w, http.StatusConflict, "signup.html", data)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// data := r.Context().Value(ctxKey).(*Data)
		h.templaterender(w, http.StatusOK, "signin.html", nil)
	case http.MethodPost:
		data := r.Context().Value(ctxKey).(*Data)
		err := r.ParseForm()
		if err != nil {
			h.errorpage(w, http.StatusInternalServerError, err)
			return
		}
		login, password := r.PostForm.Get("login"), r.PostForm.Get("password")
		data.Content = login
		if login == "" || password == "" {
			h.errorpage(w, http.StatusBadRequest, nil)
			return
		}
		user, err := h.Service.SignIn(login, password)
		if err != nil {
			data.ErrMsgs = make(map[string]string)
			switch err {
			case models.ErrNoRecord:
				data.ErrMsgs["login"] = validator.MsgUserNotFound
			case models.ErrInvalidCredentials:
				data.ErrMsgs["password"] = validator.MsgNotCorrectPassword
			default:
				h.errorpage(w, http.StatusInternalServerError, err)
				return
			}
			h.templaterender(w, http.StatusUnauthorized, "signin.html", data)
			return
		}
		cookie := &http.Cookie{
			Name:  "session",
			Value: *user.Token,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(ctxKey).(*Data)
	if data.User == (models.User{}) {
		h.errorpage(w, http.StatusUnauthorized, nil)
		return
	}
	err := h.Service.LogOut(*data.User.Token)
	if err != nil {
		h.errorpage(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
