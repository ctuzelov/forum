package service

import (
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
)

type Comments interface {
	Create(c *models.Comment) error
	Fetch(postID, userID int) ([]*models.Comment, error)
	Count(postID int) (int, error)
	GetByID(comID int) (*models.Comment, error)
	React(comID, userID int, received string) error
}

type CommentService struct {
	repo repository.Comments
}

func (s *CommentService) Create(c *models.Comment) error {
	if !(c.ParentID == 0) {
		_, err := s.repo.CommentById(c.ParentID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return models.ErrInvalidParent
			}
			return err
		}
	}
	return s.repo.InsertComment(c)
}

func (s *CommentService) GetByID(comID int) (*models.Comment, error) {
	return s.repo.CommentById(comID)
}

func (s *CommentService) Fetch(postID, userID int) ([]*models.Comment, error) {
	comments, err := s.repo.CommentsByPostId(postID)
	if err != nil {
		return comments, err
	}
	err = s.GetReplies(comments, userID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) GetReplies(comments []*models.Comment, userID int) (err error) {
	for _, comment := range comments {
		comment.Likes.Users, err = s.repo.LikesByCommentId(comment.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		comment.Likes.Count = len(comment.Likes.Users)
		comment.Dislikes.Users, err = s.repo.DislikesByCommentId(comment.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		comment.Dislikes.Count = len(comment.Dislikes.Users)
		comment.Replies, err = s.repo.RepliesByParent(comment.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		comment.UserReaction, err = s.repo.ReactionByUserId(comment.ID, userID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		err = s.GetReplies(comment.Replies, userID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
	}
	return nil
}

func (s *CommentService) Count(postID int) (int, error) {
	return s.repo.CountCommentsByPostId(postID)
}

func (s *CommentService) React(comID, userID int, received string) error {
	_, err := s.repo.CommentById(comID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return models.ErrInvalidObjectId
		}
		return err
	}
	reaction, err := s.repo.ReactionByUserId(comID, userID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}
	switch reaction {
	case "":
		err = s.repo.InsertReaction(comID, userID, received)
		if err != nil {
			return err
		}
	case received:
		err = s.repo.RemoveReaction(comID, userID)
		if err != nil {
			return err
		}
	default:
		err = s.repo.UpdateReaction(comID, userID, received)
		if err != nil {
			return err
		}
	}
	return nil
}
