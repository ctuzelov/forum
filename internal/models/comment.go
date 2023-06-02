package models

type Comment struct {
	ID       int
	PostID   int
	UserID   int
	Content  string
	ParentID int
	Creator  string
	Replies  []*Comment
	Likes    struct {
		Count int
		Users []string
	}
	Dislikes struct {
		Count int
		Users []string
	}
	UserReaction string
}
