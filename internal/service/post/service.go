package post

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/trust-me-im-an-engineer/comments/internal/cursorcode"
	"github.com/trust-me-im-an-engineer/comments/internal/domain"
	"github.com/trust-me-im-an-engineer/comments/internal/errs"
	"github.com/trust-me-im-an-engineer/comments/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func (s *Service) GetPosts(ctx context.Context, sort string, limit int32, cursor *string) (posts []*domain.Post, nextCursor string, hasNext bool, err error) {
	switch sort {
	case domain.SortOrderRating:
		if cursor != nil {
			rating, id, err := cursorcode.DecodeRatingID(*cursor)
			if err != nil {
				return nil, "", false, errs.InvalidCursor
			}

			posts, hasNext, err = s.storage.GetPostsRatingCursor(ctx, limit, rating, id)
			if err != nil {
				return nil, "", false, fmt.Errorf("storage failed to get posts by rating with cursor: %w", err)
			}
		} else {
			var err error
			posts, hasNext, err = s.storage.GetPostsRating(ctx, limit)
			if err != nil {
				return nil, "", false, fmt.Errorf("storage failed to get posts by rating: %w", err)
			}
		}
		lastPost := posts[len(posts)-1]
		nextCursor = cursorcode.EncodeRatingID(lastPost.Rating, lastPost.ID)

	case domain.SortOrderNew, domain.SortOrderOld:
		if cursor != nil {
			t, id, err := cursorcode.DecodeTimeID(*cursor)
			if err != nil {
				return nil, "", false, errs.InvalidCursor
			}

			posts, hasNext, err = s.storage.GetPostsTimeCursor(ctx, limit, t, id, true)
			if err != nil {
				return nil, "", false, fmt.Errorf("storage failed to get posts by time with cursor: %w", err)
			}
		} else {
			var err error
			posts, hasNext, err = s.storage.GetPostsTime(ctx, limit, true)
			if err != nil {
				return nil, "", false, fmt.Errorf("storage failed to get posts by time: %w", err)
			}
		}
		lastPost := posts[len(posts)-1]
		nextCursor = cursorcode.EncodeTimeID(lastPost.CreatedAt, lastPost.ID)
	}

	return posts, nextCursor, hasNext, nil
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage}
}

func (s *Service) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	post, err := s.storage.GetPost(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("storage failed to get post: %w", err)
	}
	return post, nil
}

func (s *Service) CreatePost(ctx context.Context, createPostInput *domain.CreatePostInput) (*domain.Post, error) {
	post, err := s.storage.CreatePost(ctx, createPostInput)
	if err != nil {
		return nil, fmt.Errorf("storage failed to create post: %w", err)
	}

	slog.Debug("post created", "postID", post.ID, "authorID", post.AuthorID)
	return post, nil
}

func (s *Service) UpdatePost(ctx context.Context, updatePostInput *domain.UpdatePostInput) (*domain.Post, error) {
	post, err := s.storage.UpdatePost(ctx, updatePostInput)
	if err != nil {
		return nil, fmt.Errorf("storage failed to update post: %w", err)
	}

	slog.Debug("post updated", "postID", post.ID)
	return post, nil
}

func (s *Service) DeletePost(ctx context.Context, id int) error {
	err := s.storage.DeletePost(ctx, id)
	if err != nil {
		return fmt.Errorf("storage failed to delete post: %w", err)
	}

	slog.Debug("post deleted", "postID", id)
	return nil
}

func (s *Service) SetCommentsRestricted(ctx context.Context, internalID int, restricted bool) (*domain.Post, error) {
	post, err := s.storage.SetCommentsRestricted(ctx, internalID, restricted)
	if err != nil {
		return nil, fmt.Errorf("storage failed to set comments restricted: %w", err)
	}

	slog.Debug("comments restriction changed", "postID", post.ID, "restricted", post.CommentsRestricted)
	return post, nil
}

func (s *Service) VotePost(ctx context.Context, internalInput *domain.PostVote) (*domain.Post, error) {
	post, err := s.storage.VotePost(ctx, internalInput)
	if err != nil {
		return nil, fmt.Errorf("storage failed to vote post: %w", err)
	}
	return post, nil
}
