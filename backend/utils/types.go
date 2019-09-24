package utils

import (
	"github.com/lib/pq"
	"time"
)

// TimeToNullTime converts a regular time.Time to a nullable pq.NullTime
func TimeToNullTime(t *time.Time) pq.NullTime {
	if t == nil || t.IsZero() {
		return pq.NullTime{Time: time.Time{}, Valid: false}
	}
	return pq.NullTime{Time: *t, Valid: true}
}
