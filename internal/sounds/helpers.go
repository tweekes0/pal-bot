package sounds

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// Delete a file based on the supplied filename
func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

const (
	mmssRegEx = "[0-5][0-9]:[0-5][0-9]"
	hhmmssRegEx = "0[0-9]:[0-5][0-9]:[0-5][0-9]"
)

// Hashes file and return sha256hash as a string
func HashFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", nil
	}	
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil
	}

	s := fmt.Sprintf("%x", h.Sum(nil))
	return s, nil
}

// Gets the filename of a file before the extension.
func getFilename(filepath string) string {
	fp := strings.Split(filepath, "/")
	fn := strings.Split(fp[len(fp)-1], ".")[0]

	return fn
}

// Converts the starttime string(00:00:00 or 00:00) to time.Duration 
func startTimeToDuration(t string) (time.Duration, error) {
	b1, _ := regexp.MatchString(mmssRegEx, t)
	b2, _ := regexp.MatchString(hhmmssRegEx, t)
	if !b1 && !b2 {
		return 0, ErrInvalidStartTime
	}

	ss := strings.Split(t, ":")
	if len(ss) > 2 {
		h := ss[0]
		m := ss[1]
		s := ss[2]
		ss := fmt.Sprintf("%02vh%02vm%02vs", h, m, s)
		return time.ParseDuration(ss)
	}

	m := ss[0]
	s := ss[1]

	dur := fmt.Sprintf("%02vm%02vs", m, s)
	return time.ParseDuration(dur)
}