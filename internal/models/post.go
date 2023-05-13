package models

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	Creator   string
	CmntCount int
	Likes     struct {
		Count int
		Users []string
	}
	Dislikes struct {
		Count int
		Users []string
	}
	Comments   []*Comment
	Categories []string
	CatID      []int
}

