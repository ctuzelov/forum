package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type PostRepo struct {
	*sql.DB
}

type Posts interface {
	InsertPost(p *models.Post) (int, error)
	FetchPosts() ([]*models.Post, error)
	PostById(id int) (*models.Post, error)
	InsertReaction(postID, userID int, reaction string) error
	ReactionByUserId(postID, userID int) (string, error)
	RemoveReaction(postID, userID int) error
	UpdateReaction(postID, userID int, reaction string) error
	LikesByPostId(postID int) ([]string, error)
	DislikesByPostId(postID int) ([]string, error)
	InsertCategory(postID int, catID []int) error
	CategoriesById(postID int) ([]string, error)
	Filter(catID []int) ([]*models.Post, error)
	PostsByUserId(userID int) ([]*models.Post, error)
	UserLikedPosts(userID int) ([]*models.Post, error)
	FetchCategories() ([]*string, error)
}

func (r *PostRepo) InsertPost(p *models.Post) (int, error) {
	query := `INSERT INTO posts
	VALUES(NULL, ?, ?, ?, ?)`
	res, err := r.Exec(query, p.UserID, p.Title, p.Content, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *PostRepo) PostById(id int) (*models.Post, error) {
	query := `SELECT *, (SELECT name FROM users WHERE users.id = posts.user_id) FROM posts WHERE id = ?`
	p := &models.Post{}
	err := r.QueryRow(query, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Creator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) CategoriesById(postID int) ([]string, error) {
	query := `SELECT cat_id, (
		SELECT category FROM categories WHERE categories.id = post_cat.cat_id
		)
	FROM post_cat WHERE post_id=?`
	var categories []string
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var catid int64
		var cat string
		err = rows.Scan(&catid, &cat)
		categories = append(categories, cat)
		if err != nil {
			return nil, err
		}
	}
	return categories, nil
}

func (r *PostRepo) FetchPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	query := `SELECT *, (
		SELECT name FROM users WHERE users.id = posts.user_id
		)
	FROM posts`
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Creator)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, err
}

func (r *PostRepo) InsertCategory(postID int, catID []int) error {
	query := `INSERT INTO post_cat (post_id, cat_id) VALUES(?, ?)`
	for _, i := range catID {
		_, err := r.Exec(query, postID, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) InsertReaction(postID, userID int, reaction string) error {
	query := `INSERT INTO reactions (post_id, user_id, type) VALUES(?, ?, ?)`
	_, err := r.Exec(query, postID, userID, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) ReactionByUserId(postID, userID int) (string, error) {
	query := `SELECT type FROM reactions WHERE post_id = ? AND user_id = ?`
	var reaction string
	err := r.QueryRow(query, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reaction, models.ErrNoRecord
		}
		return reaction, err
	}
	return reaction, nil
}

func (r *PostRepo) RemoveReaction(postID, userID int) error {
	query := "DELETE FROM reactions WHERE post_id = ? AND user_id = ?"
	_, err := r.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) UpdateReaction(postID, userID int, reaction string) error {
	query := `UPDATE reactions
	SET type = ?
	WHERE post_id = ? AND user_id = ?`
	_, err := r.Exec(query, reaction, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) LikesByPostId(postID int) ([]string, error) {
	query := `SELECT users.name FROM reactions
	JOIN users ON reactions.user_id=users.id
	WHERE reactions.post_id=? AND reactions.type="like"`
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	likes := []string{}
	defer rows.Close()
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		likes = append(likes, username)
	}
	return likes, err
}

func (r *PostRepo) DislikesByPostId(postID int) ([]string, error) {
	query := `SELECT users.name FROM reactions
	JOIN users ON reactions.user_id=users.id
	WHERE reactions.post_id=? AND reactions.type="dislike"`
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	dislikes := []string{}
	defer rows.Close()
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		dislikes = append(dislikes, username)
	}
	return dislikes, err
}

func (r *PostRepo) Filter(catID []int) ([]*models.Post, error) {
	newpost := []*models.Post{}
	query := `SELECT posts.id, posts.user_id, posts.title, posts.content, posts.date, users.name  FROM posts 
	JOIN post_cat ON posts.id=post_cat.post_id
	JOIN users ON users.id=posts.user_id
	WHERE post_cat.cat_id=?;`
	for _, i := range catID {
		rows, err := r.Query(query, i)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			p := &models.Post{}
			err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Creator)
			if err != nil {
				return nil, err
			}
			if newpost != nil && IsUnique(p.ID, newpost) {
				newpost = append(newpost, p)
			}
		}
	}
	return newpost, nil
}

func IsUnique(postID int, posts []*models.Post) bool {
	for _, post := range posts {
		if postID == post.ID {
			return false
		}
	}
	return true
}

func (r *PostRepo) PostsByUserId(userID int) ([]*models.Post, error) {
	query := `SELECT *, (
        SELECT name FROM users WHERE users.id = posts.user_id
        )
    FROM posts WHERE user_id = ?`
	rows, err := r.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Creator)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepo) UserLikedPosts(userID int) ([]*models.Post, error) {
	query := `SELECT posts.*, (SELECT name FROM users WHERE users.id = posts.user_id)
    FROM posts
    INNER JOIN reactions ON reactions.post_id = posts.id
    WHERE reactions.user_id = ? and reactions.type = "like";
    `
	rows, err := r.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*models.Post{}
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Creator)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepo) FetchCategories() ([]*string, error) {
	categories := []*string{}
	query := `SELECT category FROM categories`
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c *string
		err := rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
