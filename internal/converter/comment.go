package converter

import (
	"strconv"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
)

func Comment_DomainToModel(d *domain.Comment) *model.Comment {
	m := &model.Comment{
		ID:        strconv.Itoa(d.ID),
		PostID:    strconv.Itoa(d.PostID),
		AuthorID:  d.AuthorID,
		CreatedAt: d.CreatedAt,
		Rating:    d.Rating,
		ParentID:  nil,
	}

	if d.Text != nil {
		m.Text = *d.Text
	}
	return m
}

func CreateCommentInput_ModelToDomain(m *model.CreateCommentInput) *domain.CreateCommentInput {
	postID, _ := strconv.Atoi(m.PostID)
	d := &domain.CreateCommentInput{
		PostID:   postID,
		AuthorID: m.AuthorID,
		Text:     m.Text,
		ParentID: nil,
	}
	if m.ParentID != nil {
		parentID, _ := strconv.Atoi(*m.ParentID)
		d.ParentID = &parentID
	}
	return d
}

func UpdateCommentInput_ModelToDomain(m *model.UpdateCommentInput) *domain.UpdateCommentInput {
	id, _ := strconv.Atoi(m.ID) // id already validated
	return &domain.UpdateCommentInput{
		ID:   id,
		Text: m.Text,
	}
}

func ModelVoteInputToDomainCommentVote(m *model.VoteInput) *domain.CommentVote {
	id, _ := strconv.Atoi(m.ID) // id already validated
	return &domain.CommentVote{
		Vote: domain.Vote{
			ID:      id,
			VoterID: m.VoterID,
			Value:   int8(m.Value),
		},
	}
}
