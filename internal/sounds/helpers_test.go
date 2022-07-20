package sounds

import (
	"testing"
	"time"

	test "github.com/tweekes0/pal-bot/internal/testing"
)

func TestGetFilename(t *testing.T) {
	tt := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "get filename for file in path",
			input:       "/path/to/file.go",
			expected:    "file",
		},
		{
			description: "get filename for normal file name without path",
			input:       "file.go",
			expected:    "file",
		},
		{
			description: "get filename for relative file path",
			input:       "./file.go",
			expected:    "file",
		},
		{
			description: "get filename for empty file path",
			input:       "",
			expected:    "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			got := getFilename(tc.input)
			test.AssertType(t, got, tc.expected)
		})
	}
}

func startToDurationTestFunc(t *testing.T, input, expectedInput string, expectedErr error) {
	t.Parallel()

	got, err := startTimeToDuration(input)
	expected, _ := time.ParseDuration(expectedInput)

	test.AssertError(t, err, expectedErr)
	if err == nil {
		test.AssertType(t, got, expected)
	}

}

func TestStartTimeToDuration(t *testing.T) {

	tt := []struct {
		description string
		input       string
		expected    string
		expectedErr error
		testFunc    func(*testing.T, string, string, error)
	}{
		{
			description: "minute and second to duration",
			input:       "12:12",
			expected:    "12m12s",
			expectedErr: nil,
			testFunc:    startToDurationTestFunc,
		},
		{
			description: "hour, minute and second to duration",
			input:       "01:12:30",
			expected:    "01h12m30s",
			expectedErr: nil,
			testFunc:    startToDurationTestFunc,
		},
		{
			description: "invalid start time",
			input:       "99:99",
			expected:    "09m09s",
			expectedErr: ErrInvalidStartTime,
			testFunc:    startToDurationTestFunc,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			tc.testFunc(t, tc.input, tc.expected, tc.expectedErr)
		})
	}
}
