package sounds

import (
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

func getFilename(filepath string) string {
	fp := strings.Split(filepath, "/")
	fn := strings.Split(fp[len(fp)-1], ".")[0]

	return fn
}