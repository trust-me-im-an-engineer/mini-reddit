package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	AuthorID  uuid.UUID `db:"author_id"`
	Text      *string   `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	Rating    int32     `db:"rating"`
	Deleted   bool      `db:"deleted"`
	ParentID  *int      `db:"parent_id"`
}

type CreateCommentInput struct {
	PostID   int
	AuthorID uuid.UUID
	Text     string
	ParentID *int
}

type UpdateCommentInput struct {
	ID   int
	Text string
}

type CommentVote struct {
	Vote
}
