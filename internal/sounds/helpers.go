package sounds

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func deleteFile(file *os.File) error {
	err := os.Remove(file.Name())
	if err != nil {
		return err
	}

	return nil
}

// hashes file and return sha256hash as a string
func HashFile(file *os.File) (string, error) {
	f, err := os.Open(file.Name())
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

func getFilename(filepath string) string {
	fp := strings.Split(filepath, "/")
	fn := strings.Split(fp[len(fp)-1], ".")[0]

	return fn
}