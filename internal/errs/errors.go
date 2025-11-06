package errs

import (
	"errors"
	"fmt"
)

var (
	PostNotFound         = errors.New("post not found")
	CommentNotFound      = errors.New("comment not found")
	InvalidID            = errors.New("id must be valid integer")
	CommentDeleted       = errors.New("comment is deleted")
	ParentCommentDeleted = errors.New("cannot reply to deleted comment")
	InvalidCursor        = errors.New("invalid cursor")
	InternalServer       = errors.New("internal server error")
)

func InvalidInputWrap(err error) error {
	return fmt.Errorf("Invalid input: %w", err)
}
