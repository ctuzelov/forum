package service

import (
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo     repository.Posts
	cmntRepo repository.Comments
}

type Posts interface {
	GetAll() ([]*models.Post, error)
	Create(p *models.Post) (int, error)
	GetById(id int) (*models.Post, error)
	React(postID, userID int, reaction string) error
	GetFiltered(catID []int) ([]*models.Post, error)
	GetUserCreated(userID int) ([]*models.Post, error)
	GetUserLiked(userID int) ([]*models.Post, error)
	GetAllCategories() ([]*string, error)
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	posts, err := s.repo.FetchPosts()
	for _, post := range posts {
		err = s.fillpost(post)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) GetFiltered(catID []int) ([]*models.Post, error) {
	posts, err := s.repo.Filter(catID)
	for _, post := range posts {
		err = s.fillpost(post)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) Create(p *models.Post) (int, error) {
	id, err := s.repo.InsertPost(p)
	if err != nil {
		return id, err
	}
	err = s.repo.InsertCategory(id, p.CatID)
	return id, err
}

func (s *PostService) GetById(id int) (*models.Post, error) {
	post, err := s.repo.PostById(id)
	if err != nil {
		return nil, err
	}
	err = s.fillpost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) React(postID, userID int, received string) error {
	reaction, err := s.repo.ReactionByUserId(postID, userID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}
	switch reaction {
	case "":
		err = s.repo.InsertReaction(postID, userID, received)
		if err != nil {
			return err
		}
	case received:
		err = s.repo.RemoveReaction(postID, userID)
		if err != nil {
			return err
		}
	default:
		err = s.repo.UpdateReaction(postID, userID, received)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) fillpost(post *models.Post) (err error) {
	post.CmntCount, err = s.cmntRepo.CountCommentsByPostId(post.ID)
	if err != nil {
		return err
	}
	post.Likes.Users, err = s.repo.LikesByPostId(post.ID)
	if err != nil {
		return err
	}
	post.Likes.Count = len(post.Likes.Users)
	post.Dislikes.Users, err = s.repo.DislikesByPostId(post.ID)
	if err != nil {
		return err
	}
	post.Dislikes.Count = len(post.Dislikes.Users)
	post.Categories, err = s.repo.CategoriesById(post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetUserCreated(userID int) ([]*models.Post, error) {
	posts, err := s.repo.PostsByUserId(userID)
	for _, post := range posts {
		err = s.fillpost(post)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) GetUserLiked(userID int) ([]*models.Post, error) {
	posts, err := s.repo.UserLikedPosts(userID)
	for _, post := range posts {
		err = s.fillpost(post)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) GetAllCategories() ([]*string, error) {
	return s.repo.FetchCategories()
}
