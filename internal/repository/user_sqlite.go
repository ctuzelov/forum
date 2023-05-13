package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Users interface {
	InsertUser(models.User) error
	UserByEmail(email string) (models.User, error)
	UserByName(name string) (models.User, error)
	UserByToken(token string) (models.User, error)
	UserById(id int) (models.User, error)
	SetToken(id int, token string) error
	RemoveToken(token string) error
}

type UserRepo struct {
	*sql.DB
}

func (r *UserRepo) InsertUser(m models.User) error {
	query := `INSERT INTO users (email, name, password)
	VALUES(?, ?, ?);`
	if _, err := r.Exec(query, m.Email, m.Name, m.Password); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UserByEmail(email string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = email`
	var user models.User
	err := r.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token, &user.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

func (r *UserRepo) UserByName(name string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = name`
	var user models.User
	err := r.QueryRow(query, name).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token, &user.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

func (r *UserRepo) UserByToken(token string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = token AND expires > DATETIME('now')`
	var user models.User
	err := r.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token, &user.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

func (r *UserRepo) UserById(id int) (models.User, error) {
	query := `SELECT * FROM users
	WHERE id = ?`
	user := models.User{}
	err := r.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token, &user.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

func (r *UserRepo) SetToken(id int, token string) error {
	query := `UPDATE users
	SET token = ?, expires = DATETIME('now', '+8 hours')
	WHERE ? = id` // expiration datetime = now + 2 hours
	if _, err := r.Exec(query, token, id); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) RemoveToken(token string) error {
	query := `UPDATE users
	SET token = NULL
	WHERE token = ?`
	_, err := r.Exec(query, token)
	return err
}
