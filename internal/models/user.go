package models

type User struct {
	ID       int
	Email    string
	Name     string
	Password string
	ConPassword string
	Token    *string
	Expires  *string
}
