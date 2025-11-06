package validator

import (
	"errors"
	"strconv"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
)

const (
	MaxTitleLen         = 200
	MaxContentLen       = 20000
	MaxPostsLimit int32 = 100
)

var (
	EmptyTitleErr     = errors.New("post title cannot be empty")
	EmptyContentErr   = errors.New("post content cannot be empty")
	TooLongTitleErr   = errors.New("post title cannot be longer than " + strconv.Itoa(MaxTitleLen) + " characters")
	TooLongContentErr = errors.New("post content cannot be longer than " + strconv.Itoa(MaxContentLen) + " characters")
	TooBigPostLimit   = errors.New("post limit cannot be longer than " + strconv.Itoa(int(MaxPostsLimit)))
)

func ValidateCreatePostInput(in model.CreatePostInput) error {
	if err := validateTitle(in.Title); err != nil {
		return err
	}
	return validateContent(in.Content)
}

func ValidateUpdatePostInput(in model.UpdatePostInput) error {
	if _, err := strconv.Atoi(in.ID); err != nil {
		return errs.InvalidID
	}

	if in.Title == nil && in.Content == nil {
		return NothingToUpdateErr
	}

	if in.Title != nil {
		if err := validateTitle(*in.Title); err != nil {
			return err
		}
	}
	if in.Content != nil {
		if err := validateContent(*in.Content); err != nil {
			return err
		}
	}

	return nil
}

func ValidatePostsInput(sort model.SortOrder, limit int32, cursor *string) error {
	if limit < 0 {
		return NegativeLimit
	}
	if limit > MaxPostsLimit {
		return TooBigPostLimit
	}
	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return EmptyTitleErr
	}
	if len(title) > MaxTitleLen {
		return TooLongTitleErr
	}
	return nil
}

func validateContent(content string) error {
	if content == "" {
		return EmptyContentErr
	}
	if len(content) > MaxContentLen {
		return TooLongContentErr
	}
	return nil
}
