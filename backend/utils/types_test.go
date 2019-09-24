package utils

import (
	"github.com/lib/pq"
	"testing"
	"time"
)

func TestTimeToNullTime(t *testing.T) {
	input := time.Time{}
	expected := pq.NullTime{Time: time.Time{}, Valid: false}
	result := TimeToNullTime(&input)
	if result != expected {
		t.Error(
			"Expected", expected,
			"but got", result,
			"for input", input,
		)
	}

	input = time.Now()
	expected = pq.NullTime{Time: input, Valid: true}
	result = TimeToNullTime(&input)
	if result != expected {
		t.Error(
			"Expected", expected,
			"but got", result,
			"for input", input,
		)
	}
}
