package converter

import (
	"github.com/trust-me-im-an-engineer/comments/graph/model"
	"github.com/trust-me-im-an-engineer/comments/internal/domain"
)

func pageInfo_DomainToModel(d *domain.PageInfo) *model.PageInfo {
	return &model.PageInfo{
		HasNextPage: d.HasNext,
		EndCursor:   d.EndCursor,
	}
}
