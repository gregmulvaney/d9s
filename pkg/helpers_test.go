package pkg

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	test := []struct {
		input    int64
		expected string
	}{
		{500, "500B"},
		{1024, "1.00KiB"},
		{1048576, "1.00MiB"},
		{1073741824, "1.00GiB"},
	}
	for _, test := range test {
		result := FormatBytes(test.input)
		if result != test.expected {
			t.Errorf("Input: %d Expected: %s Got: %s", test.input, test.expected, result)
		}
	}
}
