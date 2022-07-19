package sounds

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/h2non/filetype"
	mp3 "github.com/hajimehoshi/go-mp3"
)

// Checks whether a file is an mp3 based on its headers
func isMP3(f *os.File) bool {
	buf, _ := ioutil.ReadFile(f.Name())
	t, _ := filetype.Audio(buf)
	
	return t.MIME.Value == "audio/mpeg" 
}

// Checks the duration of a MP3 file
func getMP3Duration(f *os.File) (int, error) {
	rd, err := os.Open(f.Name())
	if err != nil {
		return 0, nil
	}
	defer rd.Close()

	d, err  := mp3.NewDecoder(rd)
	if err != nil {
		return 0, err
	}

	samples := d.Length() / 4
	dur := samples / int64(d.SampleRate())

	return int(dur), nil
}

// Ensures the file is an MP3 and less than specified duration(in seconds)
func validateMP3(f *os.File, maxDuration int) error {
	if !isMP3(f) {
		return ErrInvalidFile
	}

	dur, err := getMP3Duration(f)
	if err != nil {
		return err
	}

	if dur >  maxDuration {
		return ErrLengthTooLong
	}

	return nil
}

func DownloadFileFromURL(name, url string, maxDuration int) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	f, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(f, resp.Body); err != nil {
		return nil, err
	}	
	defer resp.Body.Close()

	if err = validateMP3(f, maxDuration); err != nil {
		DeleteFile(f.Name())
		return nil, err
	}

	f.Seek(0, io.SeekCurrent)

	return f, nil
}

func MP3ToDCA(path string, f *os.File) (*os.File, error) {
	c1 := exec.Command("ffmpeg", "-i", f.Name(), "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	c2 := exec.Command("dca")

	var err error
	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return nil, err
	}

	name := getFilename(f.Name())
	file, err := os.Create(fmt.Sprintf("%v/%v.dca", path, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c2.Stdout = file

	err = c2.Start()
	if err != nil {
		return nil, err
	}

	err = c1.Run()
	if err != nil {
		return nil, err
	}

	err = c2.Wait()
	if err != nil {
		return nil, err
	}

	return file, err
}