package converter

import (
	"strconv"

	"github.com/trust-me-im-an-engineer/comments/graph/model"
	"github.com/trust-me-im-an-engineer/comments/internal/domain"
)

func Post_DomainToModel(d *domain.Post) *model.Post {
	return &model.Post{
		ID:                 strconv.Itoa(d.ID),
		AuthorID:           d.AuthorID,
		Title:              d.Title,
		Content:            d.Content,
		CreatedAt:          d.CreatedAt,
		Rating:             d.Rating,
		CommentsCount:      d.CommentsCount,
		CommentsRestricted: d.CommentsRestricted,
	}
}

func CreatePostInput_ModelToDomain(m *model.CreatePostInput) *domain.CreatePostInput {
	return &domain.CreatePostInput{
		AuthorID: m.AuthorID,
		Title:    m.Title,
		Content:  m.Content,
	}
}

func UpdatePost_ModelToDomain(m *model.UpdatePostInput) *domain.UpdatePostInput {
	id, _ := strconv.Atoi(m.ID) // id already validated
	return &domain.UpdatePostInput{
		ID:      id,
		Title:   m.Title,
		Content: m.Content,
	}
}

func ModelVoteInputToDomainPostVote(m *model.VoteInput) *domain.PostVote {
	id, _ := strconv.Atoi(m.ID) // id already validated
	return &domain.PostVote{
		Vote: domain.Vote{
			ID:      id,
			VoterID: m.VoterID,
			Value:   int8(m.Value),
		},
	}
}

func PostConnection_DomainToModel(d *domain.PostConnection) *model.PostConnection {
	edges := make([]*model.PostEdge, len(d.Edges))
	for i, e := range d.Edges {
		edges[i] = &model.PostEdge{
			Cursor: *e.Cursor,
			Node:   Post_DomainToModel(e.Post),
		}
	}

	return &model.PostConnection{
		Edges:    edges,
		PageInfo: pageInfo_DomainToModel(d.PageInfo),
	}
}

func PostsInput(sort model.SortOrder, limit int32, cursor *string) *domain.PostsInput {
	return &domain.PostsInput{
		Sort:   domain.SortOrder(sort),
		Limit:  limit,
		Cursor: cursor,
	}
}
