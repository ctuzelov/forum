package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Comments interface {
	InsertComment(c *models.Comment) error
	CommentById(comID int) (*models.Comment, error)
	CommentsByPostId(postID int) ([]*models.Comment, error)
	CountCommentsByPostId(postID int) (int, error)
	RepliesByParent(parentID int) ([]*models.Comment, error)
	InsertReaction(comID, userID int, reaction string) error
	ReactionByUserId(comID, userID int) (string, error)
	RemoveReaction(comID, userID int) error
	UpdateReaction(comID, userID int, reaction string) error
	LikesByCommentId(comID int) ([]string, error)
	DislikesByCommentId(comID int) ([]string, error)
}

type CommentRepo struct {
	*sql.DB
}

func (r *CommentRepo) InsertComment(c *models.Comment) error {
	query := `INSERT INTO comments
	VALUES(NULL, ?, ?, ?, ?)`
	_, err := r.Exec(query, c.PostID, c.UserID, c.Content, c.ParentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) CommentById(comID int) (*models.Comment, error) {
	query := `SELECT * FROM comments WHERE id = ?`
	comment := &models.Comment{}
	err := r.QueryRow(query, comID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.ParentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return comment, nil
}

func (r *CommentRepo) CommentsByPostId(postID int) ([]*models.Comment, error) {
	query := `SELECT *, (SELECT name FROM users WHERE comments.user_id = users.id) FROM comments
	WHERE post_id = ? AND parent_id = 0`
	comments := []*models.Comment{}
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ParentID, &c.Creator)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepo) RepliesByParent(parentID int) ([]*models.Comment, error) {
	query := `SELECT *, (SELECT name FROM users WHERE comments.user_id = users.id) FROM comments
	WHERE parent_id = ?`
	replies := []*models.Comment{}
	rows, err := r.Query(query, parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ParentID, &c.Creator)
		if err != nil {
			return nil, err
		}
		replies = append(replies, c)
	}
	return replies, nil
}

func (r *CommentRepo) CountCommentsByPostId(postID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
	err := r.QueryRow(query, postID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}
	return count, nil
}

func (r *CommentRepo) InsertReaction(comID, userID int, reaction string) error {
	query := `INSERT INTO reactions (comment_id, user_id, type) VALUES(?, ?, ?)`
	_, err := r.Exec(query, comID, userID, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) ReactionByUserId(comID, userID int) (string, error) {
	query := `SELECT type FROM reactions WHERE comment_id = ? AND user_id = ?`
	var reaction string
	err := r.QueryRow(query, comID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reaction, models.ErrNoRecord
		}
		return reaction, err
	}
	return reaction, nil
}

func (r *CommentRepo) RemoveReaction(comID, userID int) error {
	query := "DELETE FROM reactions WHERE comment_id = ? AND user_id = ?"
	_, err := r.Exec(query, comID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) UpdateReaction(comID, userID int, reaction string) error {
	query := `UPDATE reactions
	SET type = ?
	WHERE comment_id = ? AND user_id = ?`
	_, err := r.Exec(query, reaction, comID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) LikesByCommentId(comID int) ([]string, error) {
	query := `SELECT users.name FROM reactions
	JOIN users ON reactions.user_id=users.id
	WHERE reactions.comment_id=? AND reactions.type="like"`
	rows, err := r.Query(query, comID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
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

func (r *CommentRepo) DislikesByCommentId(comID int) ([]string, error) {
	query := `SELECT users.name FROM reactions
	JOIN users ON reactions.user_id=users.id
	WHERE reactions.comment_id=? AND reactions.type="dislike"`
	rows, err := r.Query(query, comID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
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
