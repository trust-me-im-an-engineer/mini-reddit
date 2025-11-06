package validator

import (
	"errors"
	"strconv"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
)

var (
	InvalidVoteValueErr = errors.New("vote value must be 1 or -1")
	NegativeLimit       = errors.New("limit must be positive")
	NothingToUpdateErr  = errors.New("at least one field needed to update")
)

func ValidateVoteInput(in model.VoteInput) error {
	if _, err := strconv.Atoi(in.ID); err != nil {
		return errs.InvalidID
	}
	if in.Value != 1 && in.Value != -1 {
		return InvalidVoteValueErr
	}
	return nil
}
