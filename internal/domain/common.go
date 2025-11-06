// package model contains internal model representation different from graph/model
package domain

import (
	"github.com/google/uuid"
)

type SortOrder string

const (
	SortOrderRating SortOrder = "RATING"
	SortOrderNew    SortOrder = "NEW"
	SortOrderOld    SortOrder = "OLD"
)

type Vote struct {
	ID      int       `db:"id"`
	VoterID uuid.UUID `db:"voter_id"`
	// +1 for upvote, -1 for downvote
	Value int8 `db:"value"`
}

type PageInfo struct {
	HasNext   bool
	EndCursor *string
}
