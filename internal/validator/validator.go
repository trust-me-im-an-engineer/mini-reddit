package validator

import (
	"errors"
	"strconv"

	"github.com/trust-me-im-an-engineer/comments/graph/model"
	"github.com/trust-me-im-an-engineer/comments/internal/errs"
)

const (
	MaxTitleLen         = 200
	MaxContentLen       = 20000
	MaxCommentLen       = 2000
	MaxPageLimit  int32 = 100
)

var (
	EmptyTitleErr       = errors.New("post title cannot be empty")
	EmptyContentErr     = errors.New("post content cannot be empty")
	TooLongTitleErr     = errors.New("post title cannot be longer than " + strconv.Itoa(MaxTitleLen) + " characters")
	TooLongContentErr   = errors.New("post content cannot be longer than " + strconv.Itoa(MaxContentLen) + " characters")
	NothingToUpdateErr  = errors.New("at least one field needed to update")
	InvalidVoteValueErr = errors.New("vote value must be 1 or -1")
	EmptyCommentErr     = errors.New("comment cannot be empty")
	TooLongCommentErr   = errors.New("comment cannot be longer than " + strconv.Itoa(MaxCommentLen) + " characters")
	NegativeLimit       = errors.New("limit must be positive")
	TooBigPostLimit     = errors.New("post limit cannot be longer than " + strconv.Itoa(int(MaxPageLimit)))
)

func ValidateCreatePostInput(in model.CreatePostInput) error {
	if err := validateTitle(in.Title); err != nil {
		return err
	}
	return validateContent(in.Content)
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

func ValidateVoteInput(in model.VoteInput) error {
	if _, err := strconv.Atoi(in.ID); err != nil {
		return errs.InvalidID
	}
	if in.Value != 1 && in.Value != -1 {
		return InvalidVoteValueErr
	}
	return nil
}

func validateCommentText(text string) error {
	if text == "" {
		return EmptyCommentErr
	}
	if len(text) > MaxContentLen {
		return TooLongCommentErr
	}
	return nil
}

func ValidateCreateCommentInput(in model.CreateCommentInput) error {
	if _, err := strconv.Atoi(in.PostID); err != nil {
		return errs.InvalidID
	}
	if err := validateCommentText(in.Text); err != nil {
		return err
	}
	if in.ParentID != nil {
		if _, err := strconv.Atoi(*in.ParentID); err != nil {
			return errs.InvalidID
		}
	}
	return nil
}

func ValidateUpdateCommentInput(in model.UpdateCommentInput) error {
	if _, err := strconv.Atoi(in.ID); err != nil {
		return errs.InvalidID
	}
	return validateCommentText(in.Text)
}

func ValidatePostsInput(sort model.SortOrder, limit int32, cursor *string) error {
	if limit < 0 {
		return NegativeLimit
	}
	if limit > MaxPageLimit {
		return TooBigPostLimit
	}
	return nil
}
