package storage

import (
	"context"

	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
)

type Storage interface {
	Post
	Comment
}

type Post interface {
	CreatePost(ctx context.Context, input *domain.CreatePostInput) (*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	UpdatePost(ctx context.Context, input *domain.UpdatePostInput) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
	SetCommentsRestricted(ctx context.Context, id int, restricted bool) (*domain.Post, error)
	VotePost(ctx context.Context, vote *domain.PostVote) (*domain.Post, error)

	GetPostsSortedByRating(ctx context.Context, limit int32, cursor *domain.PostRatingCursor) (*domain.PostsPage, error)
	GetPostsSortedByTime(ctx context.Context, limit int32, cursor *domain.PostTimeCursor, newFirst bool) (*domain.PostsPage, error)
}

type Comment interface {
	CreateComment(ctx context.Context, input *domain.CreateCommentInput) (*domain.Comment, error)
	UpdateCommentIfNotDeleted(ctx context.Context, input *domain.UpdateCommentInput) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	VoteCommentIfNotDeleted(ctx context.Context, input *domain.CommentVote) (*domain.Comment, error)
	GetComment(ctx context.Context, id int) (*domain.Comment, error)
}
