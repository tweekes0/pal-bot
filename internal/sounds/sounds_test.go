package sounds

import (
	"io/ioutil"
	"os"
	"testing"

	test "github.com/tweekes0/pal-bot/internal/testing"

	"github.com/kkdai/youtube/v2"
)

type fileTestCase struct {
	description string
	input       fileInput
	expected    error
}

type fileInput struct {
	url      string
	start    string
	duration int
}

const (
	testURL = "https://www.youtube.com/watch?v=vkFRAIKpKmE"
)

func TestDownloadYoutubeVideo(t *testing.T) {
	tt := []struct {
		description string
		input       string
		expected    error
	}{
		{
			description: "download youtube video with full url",
			input:       "https://www.youtube.com/watch?v=vkFRAIKpKmE",
			expected:    nil,
		},
		{
			description: "download youtube video id",
			input:       "vkFRAIKpKmE",
			expected:    nil,
		},
		{
			description: "download empty string",
			input:       "",
			expected:    youtube.ErrVideoIDMinLength,
		},
		{
			description: "download invalid video id",
			input:       "???????????",
			expected:    youtube.ErrInvalidCharactersInVideoID,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			f, _, got := downloadYoutubeVideo(tc.input)

			if f != nil {
				defer DeleteFile(f.Name())
			}

			test.AssertError(t, got, tc.expected)
		})
	}
}

func TestCreateAACFile(t *testing.T) {
	tt := []fileTestCase{
		{
			description: "create valid AAC file from valid youtube url",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "00:00",
				duration: 10,
			},
			expected: nil,
		},
		{
			description: "create AAC from valid youtube url and wtih no start time",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				duration: 10,
			},
			expected: ErrInvalidStartTime,
		},
		{
			description: "create AAC from valid youtube url with invalid start time",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "99:99",
				duration: 10,
			},
			expected: ErrInvalidStartTime,
		},
		{
			description: "create AAC from valid youtube url with invalid duration",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "00:00",
				duration: -1,
			},
			expected: ErrInvalidDuration,
		},

		{
			description: "create AAC from valid youtube url with duration too long",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "00:00",
				duration: 20,
			},
			expected: ErrInvalidDuration,
		},
		{
			description: "create AAC from invalid youtube url",
			input: fileInput{
				url:      "",
				start:    "00:00",
				duration: 10,
			},
			expected: youtube.ErrVideoIDMinLength,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			dir, _ := ioutil.TempDir("", "*")
			defer os.RemoveAll(dir)

			f, got := createAACFile(dir, tc.input.url, tc.input.start, tc.input.duration)
			if f != nil {
				defer DeleteFile(f.Name())
			}

			test.AssertError(t, got, tc.expected)
		})
	}
}

func TestCreateDCAFile(t *testing.T) {
	tt := []fileTestCase{
		{
			description: "create valid DCA file from valid youtube url",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "00:00",
				duration: 10,
			},
			expected: nil,
		},
		{
			description: "create valid DCA file with invalid duration",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "00:00",
				duration: -10,
			},
			expected: ErrInvalidDuration,
		},
		{
			description: "create valid DCA file with invalid starting",
			input: fileInput{
				url:      "https://www.youtube.com/watch?v=vkFRAIKpKmE",
				start:    "99:99",
				duration: 10,
			},
			expected: ErrInvalidStartTime,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			dir, _ := ioutil.TempDir("", "*")
			defer os.RemoveAll(dir)

			dca, mp3, got := CreateDCAFile(dir, tc.input.url, tc.input.start, tc.input.duration)
			if dca != nil && mp3 != nil {
				defer DeleteFile(dca.Name())
				defer DeleteFile(mp3.Name())
			}

			test.AssertError(t, got, tc.expected)
		})
	}
}

func TestCreateMP3File(t *testing.T) {
	t.Run("Create MP3 file from valid AAC", func(t *testing.T) {
		dir, _ := ioutil.TempDir("", "*")
		defer os.RemoveAll(dir)

		aac, err := createAACFile(dir, testURL, "00:00", 10)
		defer DeleteFile(aac.Name())

		test.AssertError(t, err, nil)
		if aac == nil {
			t.Fatalf("got: %v expected: %v", aac, nil)
		}

		mp3, err := createMP3File(aac)
		defer DeleteFile(mp3.Name())

		test.AssertError(t, err, nil)
	})

	t.Run("create MP3 file from invalid AAC", func(t *testing.T) {
		mp3, err := createMP3File(nil)
		test.AssertError(t, err, ErrInvalidFile)
		if mp3 != nil {
			t.Fatalf("got: %v, expected: %v", mp3, nil)
		}
	})
}

func TestLoadSound(t *testing.T) {
	t.Run("load sound for valid DCA file", func(t *testing.T) {
		dir, _ := ioutil.TempDir("", "*")
		defer os.RemoveAll(dir)

		dca, mp3, err := CreateDCAFile(dir, testURL, "00:00", 10)
		defer DeleteFile(dca.Name())
		defer DeleteFile(mp3.Name())
		test.AssertError(t, err, nil)

		_, err = LoadSound(dca.Name())
		test.AssertError(t, err, nil)
	})
}