package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID                 int       `db:"id"`
	AuthorID           uuid.UUID `db:"author_id"`
	Title              string    `db:"title"`
	Content            string    `db:"content"`
	CreatedAt          time.Time `db:"created_at"`
	Rating             int32     `db:"rating"`
	CommentsCount      int32     `db:"comments_count"`
	CommentsRestricted bool      `db:"comments_restricted"`
}

type CreatePostInput struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

type UpdatePostInput struct {
	ID      int
	Title   *string
	Content *string
}

type PostVote struct {
	Vote
}

type PostsInput struct {
	Sort   SortOrder
	Limit  int32
	Cursor *string
}

type PostEdge struct {
	Cursor *string
	Post   *Post
}

type PostConnection struct {
	Edges    []*PostEdge
	PageInfo *PageInfo
}

type PostTimeCursor struct {
	Time time.Time
	ID   int
}

type PostRatingCursor struct {
	Rating int32
	ID     int
}

type PostsPage struct {
	Posts   []*Post
	HasNext bool
}
