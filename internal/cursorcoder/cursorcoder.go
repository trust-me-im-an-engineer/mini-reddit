package cursorcoder

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/trust-me-im-an-engineer/comments/internal/domain"
)

// EncodeTimeID is a convenience helper.
func EncodeTimeID(t time.Time, id int) string {
	return encodeParts(t.Format(time.RFC3339Nano), id)
}

// DecodeTimeID decodes a time|id cursor.
func DecodeTimeID(s string) (*domain.PostTimeCursor, error) {
	parts, err := decodeParts(s)
	if err != nil {
		return nil, err
	}
	if len(parts) != 2 {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	cursor := &domain.PostTimeCursor{
		Time: t,
		ID:   id,
	}
	return cursor, nil
}

func EncodeRatingID(rating int32, id int) string {
	return encodeParts(rating, id)
}

func DecodeRatingID(s string) (*domain.PostRatingCursor, error) {
	parts, err := decodeParts(s)
	if err != nil {
		return nil, err
	}
	if len(parts) != 2 {
		return nil, err
	}

	r, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	cursor := &domain.PostRatingCursor{
		Rating: int32(r),
		ID:     id,
	}
	return cursor, nil
}

// encodeParts encodes arbitrary values as base64("val1|val2|...")
func encodeParts(values ...any) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return base64.RawStdEncoding.EncodeToString([]byte(strings.Join(parts, "|")))
}

// decodeParts decodes base64("val1|val2|...") back into []string
func decodeParts(s string) ([]string, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "|"), nil
}
