package models

type Reaction struct {
	ID        int
	PostID    *int
	CommentID *int
	UserID    int
	Type      string
}
