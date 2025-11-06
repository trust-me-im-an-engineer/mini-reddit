package validator

import (
	"errors"
	"strconv"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
)

const (
	MaxCommentLen = 2000
)

var (
	EmptyCommentErr   = errors.New("comment cannot be empty")
	TooLongCommentErr = errors.New("comment cannot be longer than " + strconv.Itoa(MaxCommentLen) + " characters")
)

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
