package repository

import "database/sql"

type Repo struct {
	Users
	Posts
	Comments
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Users:    &UserRepo{db},
		Posts:    &PostRepo{db},
		Comments: &CommentRepo{db},
	}
}
