package subscription

import (
	"log/slog"
	"sync"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph/model"
)

type Service struct {
	Subscribers map[int]map[chan *model.Comment]struct{}
	mu          sync.RWMutex
}

func NewService() *Service {
	return &Service{
		Subscribers: make(map[int]map[chan *model.Comment]struct{}),
	}
}

func (s *Service) PublishComment(postID int, comment *model.Comment) {
	s.mu.RLock()
	subsMap, exists := s.Subscribers[postID]
	if !exists || len(subsMap) == 0 {
		s.mu.RUnlock()
		return
	}

	subscribers := make([]chan *model.Comment, 0, len(subsMap))
	for ch := range subsMap {
		subscribers = append(subscribers, ch)
	}
	s.mu.RUnlock()

	for _, ch := range subscribers {
		select {
		case ch <- comment:
		default:
		}
		slog.Warn("subscriber channel full, dropping comment", "postID", postID)
	}
}

func (s *Service) SubscribeToPost(postID int, ch chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Subscribers[postID] == nil {
		s.Subscribers[postID] = make(map[chan *model.Comment]struct{})
	}
	s.Subscribers[postID][ch] = struct{}{}
}

func (s *Service) UnsubscribeFromPost(postID int, ch chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(ch)
	delete(s.Subscribers[postID], ch)

	if len(s.Subscribers[postID]) == 0 {
		delete(s.Subscribers, postID)
	}
}
