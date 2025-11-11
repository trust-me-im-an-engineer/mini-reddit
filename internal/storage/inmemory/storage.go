package inmemory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/trust-me-im-an-engineer/mini-reddit/internal/domain"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/errs"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage implements storage.Storage with an in-memory map.
// It uses a sync.RWMutex to ensure concurrent access safety.
type Storage struct {
	posts        map[int]*domain.Post
	comments     map[int]*domain.Comment
	postVotes    map[int]map[uuid.UUID]*domain.PostVote    // PostID -> VoterID -> Vote
	commentVotes map[int]map[uuid.UUID]*domain.CommentVote // CommentID -> VoterID -> Vote

	// Mutex for concurrent access
	mu sync.RWMutex
	// Simple auto-incrementing IDs
	nextPostID    int
	nextCommentID int
}

func New() *Storage {
	return &Storage{
		posts:         make(map[int]*domain.Post),
		comments:      make(map[int]*domain.Comment),
		postVotes:     make(map[int]map[uuid.UUID]*domain.PostVote),
		commentVotes:  make(map[int]map[uuid.UUID]*domain.CommentVote),
		nextPostID:    1,
		nextCommentID: 1,
	}
}

// --- Post Methods ---

func (s *Storage) CreatePost(ctx context.Context, input *domain.CreatePostInput) (*domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	post := &domain.Post{
		ID:                 s.nextPostID,
		AuthorID:           input.AuthorID,
		Title:              input.Title,
		Content:            input.Content,
		CreatedAt:          now,
		Rating:             0,
		CommentsCount:      0,
		CommentsRestricted: false,
	}
	s.posts[post.ID] = post
	s.nextPostID++
	s.postVotes[post.ID] = make(map[uuid.UUID]*domain.PostVote) // Initialize vote map
	return post, nil
}

func (s *Storage) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, errs.PostNotFound // Updated
	}
	// Return a copy to prevent external modification without lock
	postCopy := *post
	return &postCopy, nil
}

func (s *Storage) UpdatePost(ctx context.Context, input *domain.UpdatePostInput) (*domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[input.ID]
	if !ok {
		return nil, errs.PostNotFound // Updated
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}

	postCopy := *post
	return &postCopy, nil
}

func (s *Storage) DeletePost(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return errs.PostNotFound // Updated
	}
	delete(s.posts, id)
	delete(s.postVotes, id)
	// In a real system, you would also delete related comments and comment votes.
	// For simplicity, we skip comment deletion logic here.
	return nil
}

func (s *Storage) SetCommentsRestricted(ctx context.Context, id int, restricted bool) (*domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, errs.PostNotFound // Updated
	}
	post.CommentsRestricted = restricted
	postCopy := *post
	return &postCopy, nil
}

func (s *Storage) VotePost(ctx context.Context, vote *domain.PostVote) (*domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[vote.ID]
	if !ok {
		return nil, errs.PostNotFound // Updated
	}

	votesMap := s.postVotes[vote.ID]

	// Calculate rating change
	currentVote, exists := votesMap[vote.VoterID]
	ratingChange := int32(vote.Value)

	if exists {
		// If the voter is trying to submit the same vote again (e.g., upvote when already upvoted)
		if currentVote.Value == vote.Value {
			// Unvote: delete the vote and reduce rating by the value (1 or -1)
			delete(votesMap, vote.VoterID)
			ratingChange = -int32(vote.Value)
		} else {
			// Change vote: The change is new_value - current_value (e.g., -1 - (+1) = -2 or +1 - (-1) = +2)
			ratingChange = int32(vote.Value) - int32(currentVote.Value)
			currentVote.Value = vote.Value
			// Re-assign vote in map to be safe
			votesMap[vote.VoterID] = currentVote
		}
	} else {
		// New vote: add the vote
		votesMap[vote.VoterID] = vote
	}

	post.Rating += ratingChange
	postCopy := *post
	return &postCopy, nil
}

func (s *Storage) GetPostsSortedByRating(ctx context.Context, limit int32, cursor *domain.PostRatingCursor) (*domain.PostsPage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	postsSlice := make([]*domain.Post, 0, len(s.posts))
	for _, post := range s.posts {
		postsSlice = append(postsSlice, post)
	}

	// Sort by Rating (descending), then by ID (descending) as tie-breaker
	sort.Slice(postsSlice, func(i, j int) bool {
		if postsSlice[i].Rating != postsSlice[j].Rating {
			return postsSlice[i].Rating > postsSlice[j].Rating
		}
		return postsSlice[i].ID > postsSlice[j].ID
	})

	startIndex := 0
	if cursor != nil {
		// Find the index of the post that comes *after* the cursor
		for i, post := range postsSlice {
			if post.Rating == cursor.Rating && post.ID == cursor.ID {
				startIndex = i + 1
				break
			}
		}
	}

	// Apply limit
	endIndex := startIndex + int(limit)
	if endIndex > len(postsSlice) {
		endIndex = len(postsSlice)
	}

	pagePosts := postsSlice[startIndex:endIndex]
	hasNext := endIndex < len(postsSlice)

	return &domain.PostsPage{
		Posts:   pagePosts,
		HasNext: hasNext,
	}, nil
}

func (s *Storage) GetPostsSortedByTime(ctx context.Context, limit int32, cursor *domain.PostTimeCursor, newFirst bool) (*domain.PostsPage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	postsSlice := make([]*domain.Post, 0, len(s.posts))
	for _, post := range s.posts {
		postsSlice = append(postsSlice, post)
	}

	// Sort by CreatedAt, then by ID as tie-breaker
	sort.Slice(postsSlice, func(i, j int) bool {
		// For 'New' first, sort descending by time, then descending by ID.
		if newFirst {
			if !postsSlice[i].CreatedAt.Equal(postsSlice[j].CreatedAt) {
				return postsSlice[i].CreatedAt.After(postsSlice[j].CreatedAt) // Newest first
			}
			return postsSlice[i].ID > postsSlice[j].ID // Tie-breaker
		}
		// For 'Old' first, sort ascending by time, then ascending by ID.
		if !postsSlice[i].CreatedAt.Equal(postsSlice[j].CreatedAt) {
			return postsSlice[i].CreatedAt.Before(postsSlice[j].CreatedAt) // Oldest first
		}
		return postsSlice[i].ID < postsSlice[j].ID // Tie-breaker
	})

	startIndex := 0
	if cursor != nil {
		// Find the index of the post that comes *after* the cursor
		for i, post := range postsSlice {
			if post.CreatedAt.Equal(cursor.Time) && post.ID == cursor.ID {
				startIndex = i + 1
				break
			}
		}
	}

	// Apply limit
	endIndex := startIndex + int(limit)
	if endIndex > len(postsSlice) {
		endIndex = len(postsSlice)
	}

	pagePosts := postsSlice[startIndex:endIndex]
	hasNext := endIndex < len(postsSlice)

	return &domain.PostsPage{
		Posts:   pagePosts,
		HasNext: hasNext,
	}, nil
}

// --- Comment Methods ---

func (s *Storage) CreateComment(ctx context.Context, input *domain.CreateCommentInput) (*domain.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[input.PostID]
	if !ok {
		return nil, errs.PostNotFound // Updated
	}
	if post.CommentsRestricted {
		return nil, fmt.Errorf("comments are restricted for post %d", input.PostID)
	}

	// Check parent comment if ParentID is set
	if input.ParentID != nil {
		parentComment, ok := s.comments[*input.ParentID]
		if !ok {
			return nil, errs.CommentNotFound // Parent comment not found (uses the generic CommentNotFound)
		}
		if parentComment.Text == nil {
			return nil, errs.ReplyToDeletedComment // Updated
		}
	}

	now := time.Now().UTC()
	text := input.Text
	comment := &domain.Comment{
		ID:        s.nextCommentID,
		PostID:    input.PostID,
		AuthorID:  input.AuthorID,
		Text:      &text,
		CreatedAt: now,
		Rating:    0,
		ParentID:  input.ParentID,
	}

	s.comments[comment.ID] = comment
	s.nextCommentID++
	s.commentVotes[comment.ID] = make(map[uuid.UUID]*domain.CommentVote) // Initialize vote map
	post.CommentsCount++

	return comment, nil
}

func (s *Storage) UpdateCommentIfNotDeleted(ctx context.Context, input *domain.UpdateCommentInput) (*domain.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment, ok := s.comments[input.ID]
	if !ok {
		return nil, errs.CommentNotFound // Updated
	}
	if comment.Text == nil {
		return nil, errs.CommentDeleted // Updated
	}

	newText := input.Text
	comment.Text = &newText
	commentCopy := *comment
	return &commentCopy, nil
}

func (s *Storage) DeleteComment(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment, ok := s.comments[id]
	if !ok {
		return errs.CommentNotFound // Updated
	}

	if comment.Text != nil {
		// "Delete" by setting Text to nil and removing content.
		comment.Text = nil
		// Recalculate post comments count (assuming deleted comments don't count)
		if post, ok := s.posts[comment.PostID]; ok {
			post.CommentsCount--
		}
	}

	// In a real implementation, you might not delete comment votes,
	// but keep them for historical purposes. For simplicity, we keep them here.
	return nil
}

func (s *Storage) VoteCommentIfNotDeleted(ctx context.Context, vote *domain.CommentVote) (*domain.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment, ok := s.comments[vote.ID]
	if !ok {
		return nil, errs.CommentNotFound // Updated
	}
	if comment.Text == nil {
		return nil, errs.CommentDeleted // Updated
	}

	votesMap := s.commentVotes[vote.ID]

	// Calculate rating change, similar logic to VotePost
	currentVote, exists := votesMap[vote.VoterID]
	ratingChange := int32(vote.Value)

	if exists {
		if currentVote.Value == vote.Value {
			// Unvote
			delete(votesMap, vote.VoterID)
			ratingChange = -int32(vote.Value)
		} else {
			// Change vote
			ratingChange = int32(vote.Value) - int32(currentVote.Value)
			currentVote.Value = vote.Value
			votesMap[vote.VoterID] = currentVote
		}
	} else {
		// New vote
		votesMap[vote.VoterID] = vote
	}

	comment.Rating += ratingChange
	commentCopy := *comment
	return &commentCopy, nil
}

func (s *Storage) GetComment(ctx context.Context, id int) (*domain.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comment, ok := s.comments[id]
	if !ok {
		return nil, errs.CommentNotFound // Updated
	}
	commentCopy := *comment
	return &commentCopy, nil
}

func (s *Storage) Close() {}
