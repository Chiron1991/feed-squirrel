package utils

import (
	"os"
	"testing"
)

func TestGetenv(t *testing.T) {
	envKey := "SQUIRREL_FOO"

	testData := [][3]string{
		// set, fallback, expected
		{"", "abc", "abc"},
		{"abc", "abc", "abc"},
		{"abc", "def", "abc"},
		{"123", "def", "123"},
	}

	for i := range testData {
		set := testData[i][0]
		fallback := testData[i][1]
		expected := testData[i][2]

		_ = os.Setenv(envKey, set)
		result := Getenv(envKey, fallback)
		if result != testData[i][2] {
			t.Error(
				"Expected", expected,
				"but got", result,
				"for env value", set,
				"and fallback", fallback,
			)
		}
	}
}
