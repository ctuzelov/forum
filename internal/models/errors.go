package models

import "errors"

var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateName      = errors.New("duplicate name")
	ErrInvalidParent      = errors.New("invalid parent")
)
