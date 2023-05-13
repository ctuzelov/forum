package service

import (
	"forum/internal/repository"
)

type Service struct {
	Users
	Posts
	Comments
}

func New(r *repository.Repo) *Service {
	return &Service{
		Users:    &UserService{r.Users},
		Posts:    &PostService{r.Posts, r.Comments},
		Comments: &CommentService{r.Comments},
	}
}
