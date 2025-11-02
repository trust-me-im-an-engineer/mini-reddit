package graph

import (
	"github.com/trust-me-im-an-engineer/comments/internal/service/comment"
	"github.com/trust-me-im-an-engineer/comments/internal/service/post"
	subsctiption "github.com/trust-me-im-an-engineer/comments/internal/service/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postService         *post.Service
	commentService      *comment.Service
	subscriptionService *subsctiption.Service
}

func NewResolver(post *post.Service, comment *comment.Service, subsctiption *subsctiption.Service) *Resolver {
	return &Resolver{
		postService:         post,
		commentService:      comment,
		subscriptionService: subsctiption,
	}
}
