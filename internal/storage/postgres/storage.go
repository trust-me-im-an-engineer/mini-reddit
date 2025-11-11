package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/trust-me-im-an-engineer/mini-reddit/internal/config"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage"
)

const (
	notFound            = "02000"
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"

	commentsRestricted = "90001"
	replyToDeleted     = "90002"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Storage, error) {
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

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	if s.pool != nil {
		slog.Info("closing postgres pool connection...")
		s.pool.Close()
		slog.Info("postgres pool connection closed")
	}
}

func (s *Storage) CreateComment(ctx context.Context, input *domain.CreateCommentInput) (*domain.Comment, error) {
	q := `INSERT INTO comments (post_id, author_id, text, parent_id)  
		  VALUES ($1, $2, $3, $4) RETURNING *`
	rows, _ := s.pool.Query(ctx, q, input.PostID, input.AuthorID, input.Text, input.ParentID)
	comment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Comment])
	if err != nil {
		var pgxError *pgconn.PgError
		if !errors.As(err, &pgxError) {
			return nil, err
		}

		switch pgxError.Code {
		case foreignKeyViolation:
			switch pgxError.ConstraintName {
			case "comments_post_id_fkey":
				return nil, errs.PostNotFound
			case "comments_parent_id_fkey":
				return nil, errs.ParentCommentNotFound
			}
		case commentsRestricted:
			return nil, errs.CommentsRestricted
		case replyToDeleted:
			return nil, errs.ReplyToDeletedComment
		}
	}

	return &comment, nil
}

func (s *Storage) CreatePost(ctx context.Context, input *domain.CreatePostInput) (*domain.Post, error) {
	q := `INSERT INTO posts (author_id, title, content) 
		  VALUES ($1, $2, $3) RETURNING *`
	rows, _ := s.pool.Query(ctx, q, input.AuthorID, input.Title, input.Content)
	post, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Post])
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *Storage) DeleteComment(ctx context.Context, id int) error {
	q := `UPDATE comments 
		  SET deleted = TRUE
		  WHERE id = $1`
	commandTag, err := s.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errs.CommentNotFound
	}
	return nil
}

func (s *Storage) DeletePost(ctx context.Context, id int) error {
	q := `DELETE FROM posts 
		  WHERE id = $1`
	commandTag, err := s.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errs.PostNotFound
	}
	return nil
}

func (s *Storage) GetComment(ctx context.Context, id int) (*domain.Comment, error) {
	q := `SELECT * FROM comments
		  WHERE id = $1`
	rows, _ := s.pool.Query(ctx, q, id)
	comment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Comment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.CommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

// GetPost implements storage.Storage.
func (s *Storage) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	panic("unimplemented")
}

// GetPostsSortedByRating implements storage.Storage.
func (s *Storage) GetPostsSortedByRating(ctx context.Context, limit int32, cursor *domain.PostRatingCursor) (*domain.PostsPage, error) {
	panic("unimplemented")
}

// GetPostsSortedByTime implements storage.Storage.
func (s *Storage) GetPostsSortedByTime(ctx context.Context, limit int32, cursor *domain.PostTimeCursor, newFirst bool) (*domain.PostsPage, error) {
	panic("unimplemented")
}

// SetCommentsRestricted implements storage.Storage.
func (s *Storage) SetCommentsRestricted(ctx context.Context, id int, restricted bool) (*domain.Post, error) {
	panic("unimplemented")
}

// UpdateCommentIfNotDeleted implements storage.Storage.
func (s *Storage) UpdateCommentIfNotDeleted(ctx context.Context, input *domain.UpdateCommentInput) (*domain.Comment, error) {
	panic("unimplemented")
}

// UpdatePost implements storage.Storage.
func (s *Storage) UpdatePost(ctx context.Context, input *domain.UpdatePostInput) (*domain.Post, error) {
	panic("unimplemented")
}

// VoteCommentIfNotDeleted implements storage.Storage.
func (s *Storage) VoteCommentIfNotDeleted(ctx context.Context, input *domain.CommentVote) (*domain.Comment, error) {
	panic("unimplemented")
}

// VotePost implements storage.Storage.
func (s *Storage) VotePost(ctx context.Context, vote *domain.PostVote) (*domain.Post, error) {
	panic("unimplemented")
}
