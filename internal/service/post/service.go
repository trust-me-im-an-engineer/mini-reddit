package post

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/trust-me-im-an-engineer/mini-reddit/internal/cursorcoder"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func (s *Service) GetPosts(ctx context.Context, q *domain.PostsInput) (*domain.PostConnection, error) {
	edges := make([]*domain.PostEdge, 0, q.Limit)
	var postsPage *domain.PostsPage

	switch q.Sort {
	case domain.SortOrderRating:
		var cursor *domain.PostRatingCursor
		if q.Cursor != nil {
			c, err := cursorcoder.DecodeRatingID(*q.Cursor)
			if err != nil {
				return nil, errs.InvalidCursor
			}
			cursor = c
		}

		pp, err := s.storage.GetPostsSortedByRating(ctx, q.Limit, cursor)
		if err != nil {
			return nil, fmt.Errorf("storage failed to get posts sorted by rating: %w", err)
		}

		for _, p := range pp.Posts {
			cursor := cursorcoder.EncodeRatingID(p.Rating, p.ID)
			edge := &domain.PostEdge{
				Cursor: &cursor,
				Post:   p,
			}
			edges = append(edges, edge)
		}

		postsPage = pp

	case domain.SortOrderNew, domain.SortOrderOld:
		var cursor *domain.PostTimeCursor
		if q.Cursor != nil {
			c, err := cursorcoder.DecodeTimeID(*q.Cursor)
			if err != nil {
				return nil, errs.InvalidCursor
			}
			cursor = c
		}

		newFirst := q.Sort == domain.SortOrderNew
		pp, err := s.storage.GetPostsSortedByTime(ctx, q.Limit, cursor, newFirst)
		if err != nil {
			return nil, fmt.Errorf("storage failed to get posts sorted by time: %w", err)
		}

		for _, p := range pp.Posts {
			cursor := cursorcoder.EncodeTimeID(p.CreatedAt, p.ID)
			edge := &domain.PostEdge{
				Cursor: &cursor,
				Post:   p,
			}
			edges = append(edges, edge)
		}

		postsPage = pp
	}

	connection := &domain.PostConnection{
		Edges: edges,
		PageInfo: &domain.PageInfo{
			HasNext:   postsPage.HasNext,
			EndCursor: nil,
		},
	}

	if len(edges) == 0 {
		return connection, nil
	}

	last := postsPage.Posts[len(postsPage.Posts)-1]
	var endCursor string
	if q.Sort == domain.SortOrderRating {
		endCursor = cursorcoder.EncodeRatingID(last.Rating, last.ID)
	} else {
		endCursor = cursorcoder.EncodeTimeID(last.CreatedAt, last.ID)
	}
	connection.PageInfo.EndCursor = &endCursor

	return connection, nil
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
