package service

import (
	"errors"
	"fmt"

	"forum/internal/hasher"
	"forum/internal/models"
	"forum/internal/repository"
)

type Users interface {
	SignUp(models.User) error
	SignIn(login, password string) (models.User, error)
	GetByToken(token string) (models.User, error)
	LogOut(token string) error
}

type UserService struct {
	repo repository.Users
}

func (u *UserService) SignUp(m models.User) error {
	user, err := u.repo.UserByEmail(m.Email)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return fmt.Errorf("userservice - signup #1: %w", err)
	}
	if user.Email == m.Email {
		return models.ErrDuplicateEmail
	}
	user, err = u.repo.UserByName(m.Name)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}
	if user.Name == m.Name {
		return models.ErrDuplicateName
	}
	m.Password, err = hasher.Encrypt(m.Password)
	if err != nil {
		return err
	}
	u.repo.InsertUser(m)
	return nil
}

func (s *UserService) SignIn(login, password string) (models.User, error) {
	m, err := s.repo.UserByEmail(login)
	switch {
	case errors.Is(err, models.ErrNoRecord):
		m, err = s.repo.UserByName(login)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return m, err
			}
			return m, err
		}
	case err != nil:
		return m, err
	}

	if !hasher.CorrectPassword(m.Password, password) {
		return m, models.ErrInvalidCredentials
	}

	m.Token, err = hasher.GenerateToken()
	if err != nil {
		return m, err
	}
	err = s.repo.SetToken(m.ID, *m.Token)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (s *UserService) GetByToken(token string) (models.User, error) {
	return s.repo.UserByToken(token)
}

func (s *UserService) LogOut(token string) error {
	return s.repo.RemoveToken(token)
}
