package storage

import (
	"context"
	"time"

	"github.com/trust-me-im-an-engineer/comments/internal/domain"
)

type Storage interface {
	CreatePost(ctx context.Context, input *domain.CreatePostInput) (*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
	UpdatePost(ctx context.Context, input *domain.UpdatePostInput) (*domain.Post, error)
	DeletePost(ctx context.Context, id int) error
	SetCommentsRestricted(ctx context.Context, id int, restricted bool) (*domain.Post, error)
	VotePost(ctx context.Context, vote *domain.PostVote) (*domain.Post, error)
	CreateComment(ctx context.Context, input *domain.CreateCommentInput) (*domain.Comment, error)
	UpdateCommentIfNotDeleted(ctx context.Context, input *domain.UpdateCommentInput) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	VoteCommentIfNotDeleted(ctx context.Context, input *domain.CommentVote) (*domain.Comment, error)
	GetComment(ctx context.Context, id int) (*domain.Comment, error)
	GetPostsRating(ctx context.Context, limit int32) (posts []*domain.Post, hasNext bool, err error)
	GetPostsRatingCursor(ctx context.Context, limit, rating int32, id int) (posts []*domain.Post, hasNext bool, err error)
	GetPostsTime(ctx context.Context, limit int32, newFirst bool) (posts []*domain.Post, hasNext bool, err error)
	GetPostsTimeCursor(ctx context.Context, limit int32, t time.Time, id int, newFirst bool) (posts []*domain.Post, hasNext bool, err error)
}
