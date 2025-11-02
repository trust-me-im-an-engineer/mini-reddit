package cursorcode

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EncodeTimeID is a convenience helper.
func EncodeTimeID(t time.Time, id int) string {
	return encodeParts(t.Format(time.RFC3339Nano), id)
}

// DecodeTimeID decodes a time|id cursor.
func DecodeTimeID(s string) (time.Time, int, error) {
	parts, err := decodeParts(s)
	if err != nil {
		return time.Time{}, 0, err
	}
	if len(parts) != 2 {
		return time.Time{}, 0, err
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, 0, err
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, 0, err
	}
	return t, id, nil
}

func EncodeRatingID(rating int32, id int) string {
	return encodeParts(rating, id)
}

func DecodeRatingID(s string) (rating int32, id int, err error) {
	parts, err := decodeParts(s)
	if err != nil {
		return 0, 0, err
	}
	if len(parts) != 2 {
		return 0, 0, err
	}

	r, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	id, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return int32(r), id, nil
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
