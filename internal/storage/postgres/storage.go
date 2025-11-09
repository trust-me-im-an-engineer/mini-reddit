package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/config"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage"
)

var _ storage.Storage = (*Postgres)(nil)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("connected to postgres database")

	return &Postgres{pool: pool}, nil
}

// CreateComment implements storage.Storage.
func (p *Postgres) CreateComment(ctx context.Context, input *domain.CreateCommentInput) (*domain.Comment, error) {
	panic("unimplemented")
}

// CreatePost implements storage.Storage.
func (p *Postgres) CreatePost(ctx context.Context, input *domain.CreatePostInput) (*domain.Post, error) {
	panic("unimplemented")
}

// DeleteComment implements storage.Storage.
func (p *Postgres) DeleteComment(ctx context.Context, id int) error {
	panic("unimplemented")
}

// DeletePost implements storage.Storage.
func (p *Postgres) DeletePost(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetComment implements storage.Storage.
func (p *Postgres) GetComment(ctx context.Context, id int) (*domain.Comment, error) {
	panic("unimplemented")
}

// GetPost implements storage.Storage.
func (p *Postgres) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	panic("unimplemented")
}

// GetPostsSortedByRating implements storage.Storage.
func (p *Postgres) GetPostsSortedByRating(ctx context.Context, limit int32, cursor *domain.PostRatingCursor) (*domain.PostsPage, error) {
	panic("unimplemented")
}

// GetPostsSortedByTime implements storage.Storage.
func (p *Postgres) GetPostsSortedByTime(ctx context.Context, limit int32, cursor *domain.PostTimeCursor, newFirst bool) (*domain.PostsPage, error) {
	panic("unimplemented")
}

// SetCommentsRestricted implements storage.Storage.
func (p *Postgres) SetCommentsRestricted(ctx context.Context, id int, restricted bool) (*domain.Post, error) {
	panic("unimplemented")
}

// UpdateCommentIfNotDeleted implements storage.Storage.
func (p *Postgres) UpdateCommentIfNotDeleted(ctx context.Context, input *domain.UpdateCommentInput) (*domain.Comment, error) {
	panic("unimplemented")
}

// UpdatePost implements storage.Storage.
func (p *Postgres) UpdatePost(ctx context.Context, input *domain.UpdatePostInput) (*domain.Post, error) {
	panic("unimplemented")
}

// VoteCommentIfNotDeleted implements storage.Storage.
func (p *Postgres) VoteCommentIfNotDeleted(ctx context.Context, input *domain.CommentVote) (*domain.Comment, error) {
	panic("unimplemented")
}

// VotePost implements storage.Storage.
func (p *Postgres) VotePost(ctx context.Context, vote *domain.PostVote) (*domain.Post, error) {
	panic("unimplemented")
}
