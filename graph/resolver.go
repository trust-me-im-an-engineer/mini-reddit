package graph

import (
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/comment"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/post"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postService         *post.Service
	commentService      *comment.Service
	subscriptionService *subscription.Service
}

func NewResolver(post *post.Service, comment *comment.Service, subscription *subscription.Service) *Resolver {
	return &Resolver{
		postService:         post,
		commentService:      comment,
		subscriptionService: subscription,
	}
}
