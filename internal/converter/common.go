package converter

import (
	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
)

func pageInfo_DomainToModel(d *domain.PageInfo) *model.PageInfo {
	return &model.PageInfo{
		HasNextPage: d.HasNext,
		EndCursor:   d.EndCursor,
	}
}
