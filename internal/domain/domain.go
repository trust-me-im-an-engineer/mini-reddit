// package model contains internal model representation different from graph/model
package domain

import (
	"time"

	"github.com/google/uuid"
)

type SortOrder string

const (
	SortOrderRating SortOrder = "RATING"
	SortOrderNew    SortOrder = "NEW"
	SortOrderOld    SortOrder = "OLD"
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

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	AuthorID  uuid.UUID `db:"author_id"`
	Text      *string   `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	Rating    int32     `db:"rating"`
	ParentID  *int      `db:"parent_id"`
}

type PostVote struct {
	PostID int `db:"post_id"`
	Vote
}

type CommentVote struct {
	CommentID int `db:"comment_id"`
	Vote
}

type Vote struct {
	VoterID uuid.UUID `db:"voter_id"`
	// +1 for upvote, -1 for downvote
	Value int8 `db:"value"`
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

type GetPostsInput struct {
	Sort   SortOrder
	Limit  int32
	Cursor *string
}